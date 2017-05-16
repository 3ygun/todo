package main

import (
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	db = CreateDatabase()
	defer db.Close()

	CreateStartData()
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
