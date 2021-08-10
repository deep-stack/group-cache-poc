package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	db "group.cache.poc/database"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go runHttpServer(":8080")
	go runHttpServer(":8081")
	wg.Wait()
}

func runHttpServer(addr string) {
	pool := &db.Database{}
	pool.Init()
	mux := mux.NewRouter()
	mux.HandleFunc("/get/{key}", get(pool)).Methods("GET")
	// TODO : Use query params or multipart form instead of this
	mux.HandleFunc("/set/{key}/{value}", set(pool)).Methods("POST")
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

// HTTP Middleware functions
// TODO : Return 404 in case key not found
func get(pool *db.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		value := []byte(pool.Get(mux.Vars(r)["key"]))
		if _, err := w.Write(value); err != nil {
			log.Println("HTTP Response write failed", err)
		}
	}
}

// TODO : Return Internal server error in case upsert fails
// TODO : Use query params or multipart form instead of this
func set(pool *db.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pool.Set(mux.Vars(r)["key"], mux.Vars(r)["value"])
		w.Write([]byte("KV pair upserted"))
	}
}
