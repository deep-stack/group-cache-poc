package main

import (
	_ "github.com/lib/pq"
	db "group.cache.poc/database"
)

func main() {
	test := db.Database{}
	test.Init()
}
