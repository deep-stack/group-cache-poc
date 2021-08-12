package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	gc "github.com/golang/groupcache"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	db "group.cache.poc/database"
)

var db_pool *db.Database

func main() {
	var mux = mux.NewRouter()
	addr := flag.String("addr", ":8080", "server address")
	peers := flag.String("pool", "http://localhost:8080", "server pool list (comma separated)")
	flag.Parse()

	db_pool = &db.Database{}
	db_pool.Init()
	p := strings.Split(*peers, ",")

	wg := sync.WaitGroup{}
	wg.Add(1)
	pool := gc.NewHTTPPoolOpts(p[0], nil)
	pool.Set(p...)
	mux.Handle("/_groupcache/", pool)
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
		if err := dest.SetBytes([]byte(db_pool.Get(key))); err != nil {
			log.Println("Cache Filling error :", err)
			return err
		}
		return nil
	},
))

// HTTP Middleware functions
// TODO : Return 404 in case key not found
func get(w http.ResponseWriter, r *http.Request) {
	value := []byte("")
	if err := Group.Get(context.TODO(), mux.Vars(r)["key"], gc.AllocatingByteSliceSink(&value)); err != nil {
		log.Println("Groupcache Get failed", err)
	}
	if _, err := w.Write(value); err != nil {
		log.Println("HTTP Response write failed", err)
	}
}

// TODO : Return Internal server error in case upsert fails
// TODO : Use query params or multipart form instead of this
func set(w http.ResponseWriter, r *http.Request) {
	db_pool.Set(mux.Vars(r)["key"], mux.Vars(r)["value"])
	w.Write([]byte("KV pair upserted"))
}
