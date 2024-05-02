package main

import (
	// "database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	handler := http.HandlerFunc(Server)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
