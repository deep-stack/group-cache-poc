package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	gc "github.com/mailgun/groupcache/v2"
	db "group.cache.poc/database"
)

var dbPool *db.Database

func main() {
	var mux = mux.NewRouter()
	addr := flag.String("addr", ":8080", "server address")
	peers := flag.String("pool", "http://localhost:8080", "server pool list (comma separated)")
	flag.Parse()

	dbPool = &db.Database{}
	dbPool.Init()
	p := strings.Split(*peers, ",")

	wg := sync.WaitGroup{}
	wg.Add(1)
	gcPool := gc.NewHTTPPoolOpts(p[0], &gc.HTTPPoolOptions{})
	gcPool.Set(p...)
	mux.Handle("/_groupcache/", gcPool)
	go runHttpServer(*addr, mux)
	wg.Wait()
}

func runHttpServer(addr string, mux *mux.Router) {
	mux.HandleFunc("/get/{key}", get).Methods("GET")
	// TODO : Use query params or multipart form instead of this
	mux.HandleFunc("/set/{key}/{value}", set).Methods("POST")
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}
	log.Println("Starting server on port", addr)
	log.Fatalln(srv.ListenAndServe())
}

// 64 << 20 Bytes ~ 64 MB
var Group = gc.NewGroup("foobar", 64<<20, gc.GetterFunc(
	func(ctx context.Context, key string, dest gc.Sink) error {

		log.Println("Groupcache queried for key =", key)
		if err := dest.SetBytes([]byte(dbPool.Get(key)), time.Time{}); err != nil {
			log.Println("Cache Filling error :", err)
			return err
		}
		return nil
	},
))

// HTTP Middleware functions
// TODO : Return 404 in case key not found
func get(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := []byte("")
	if err := Group.Get(context.TODO(), key, gc.AllocatingByteSliceSink(&value)); err != nil {
		log.Println("Groupcache Get failed", err)
	}
	if _, err := w.Write(value); err != nil {
		log.Println("HTTP Response write failed", err)
	}
}

// TODO : Return Internal server error in case upsert fails
// TODO : Use query params or multipart form instead of this
// TODO : Make Group.Remove to work
func set(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := mux.Vars(r)["value"]
	/*
		if err := Group.Remove(context.TODO(), key); err != nil {
			log.Println("Groupcache Remove failed", err)
		}
	*/
	dbPool.Set(key, value)
	w.Write([]byte("KV pair upserted"))
}
