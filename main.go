package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	gc "github.com/mailgun/groupcache/v2"
	db "group.cache.poc/database"
)

var dbPool *db.Database

func main() {
	addr := flag.String("addr", ":8080", "server address")
	peers := flag.String("pool", "http://localhost:8080", "server pool list (comma separated)")
	flag.Parse()

	dbPool = &db.Database{}
	dbPool.Init()
	p := strings.Split(*peers, ",")

	wg := sync.WaitGroup{}
	wg.Add(1)
	gcPool := gc.NewHTTPPool(p[0])
	gcPool.Set(p...)
	go runHttpServer(*addr)
	wg.Wait()
}

func runHttpServer(addr string) {
	http.HandleFunc("/Get", get)
	http.HandleFunc("/Set", set)
	log.Println("Starting server on port", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}

// 64 << 20 Bytes ~ 64 MB
var Group = gc.NewGroup("foobar", 64<<20, gc.GetterFunc(
	func(ctx context.Context, key string, dest gc.Sink) error {

		log.Println("Cache Miss, hitting DB for key =", key)
		if err := dest.SetBytes([]byte(dbPool.Get(key)), time.Time{}); err != nil {
			log.Println("Cache Filling error :", err)
			return err
		}
		return nil
	},
))

// TODO: Return StatusNotFound in case Key not found in psql DB
func get(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	var b []byte
	if len(key) == 0 || r.Method != "GET" {
		http.Error(w, "Invalid Method", http.StatusBadRequest)
		return
	}
	err := Group.Get(context.TODO(), key, gc.AllocatingByteSliceSink(&b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Write(b)
}

func set(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	value := query.Get("value")
	if len(key) == 0 || len(value) == 0 || r.Method != "POST" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := Group.Remove(context.TODO(), key); err != nil {
		log.Println("Groupcache Remove failed", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	dbPool.Set(key, value)
	w.Write([]byte("Key-Value pair upserted"))
}
