package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_, err := createDb()
	if err != nil {
		os.Exit(1)
	}
	err = initializeDb()
	if err != nil {
		os.Exit(1)
	}
	router := http.NewServeMux()

	router.HandleFunc("/v1/urlinfo/", Server)
	log.Fatal(http.ListenAndServe(":8080", router))

}
