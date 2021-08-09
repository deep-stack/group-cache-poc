package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

const (
	createQuery = `CREATE TABLE IF NOT EXISTS KV (
					key varchar(64),
					value varchar(64)
				   )`

	insertQuery = `INSERT INTO KV(key, value) VALUES($1, $2)`

	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	KvLength = 64
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

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Hostname, db.Port, db.User, db.Password, db.Name)
	if db.pool, err = sql.Open("postgres", psqlconn); err != nil {
		log.Fatalln(err)
	}

	if err = db.pool.Ping(); err != nil {
		log.Fatalln(err)
	}
	db.seed(10)

}

// Seeds the db with n random Key-Value pairs
func (db *Database) seed(n int) {

	rand.Seed(time.Now().UnixNano())
	if _, err := db.pool.Query(createQuery); err != nil {
		log.Fatalln(err)
	}
	var (
		stmt *sql.Stmt
		err  error
	)
	if stmt, err = db.pool.Prepare(insertQuery); err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < n; i++ {
		stmt.Exec(randString(KvLength), randString(KvLength))
	}

}

// Function to return random string of length n
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
