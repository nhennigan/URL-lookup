package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	createDb()
	initializeDb()

	router := http.NewServeMux()

	router.HandleFunc("/v1/urlinfo/{url}", Server)
	log.Fatal(http.ListenAndServe(":8080", router))

}
