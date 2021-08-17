package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

const (
	createQuery = `
				CREATE TABLE IF NOT EXISTS KV (
				key varchar(64) PRIMARY KEY,
				value varchar(64))`

	upsertQuery = `INSERT INTO KV(key, value) VALUES($1, $2)
				   ON CONFLICT(key) do 
				   UPDATE SET value = EXCLUDED.value`

	getQuery = `SELECT value from KV where key = $1`
)

// Why a Db struct : www.alexedwards.net/blog/organising-database-access
type Database struct {
	pool     *sql.DB
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Init creates connection using config.json and pings postgres db
func (db *Database) Init() {
	var (
		byt []byte
		err error
	)

	if byt, err = ioutil.ReadFile("env/config.json"); err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(byt, db)

	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Hostname,
		db.Port,
		db.User,
		db.Password,
		db.Name,
	)
	if db.pool, err = sql.Open("postgres", psqlconn); err != nil {
		log.Fatalln(err)
	}

	if err = db.pool.Ping(); err != nil {
		log.Fatalln(err)
	}

	if _, err := db.pool.Query(createQuery); err != nil {
		log.Fatalln("Failed to create table", err)
	}
}

func (db *Database) Get(key string) (value string) {
	if err := db.pool.QueryRow(getQuery, key).Scan(&value); err != nil {
		log.Println("psql db hit :", key, "not found in DB", err)
		return ""
	}

	log.Println("psql db hit -", key, ":", value, "found in DB")
	return value
}

func (db *Database) Set(key, value string) {
	if _, err := db.pool.Exec(upsertQuery, key, value); err != nil {
		log.Println("psql db - insert failed", err)
		return
	}

	log.Println("upserted to psql db -", key, ":", value)
}
