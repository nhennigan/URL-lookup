package main

//will need a HTTP handler
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

import (
	"fmt"
	"log"
	"net/http"
)

// will need a HTTP handler
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Check URL")
}

// will server some sort of webpage
func main() {
	http.HandleFunc("/v1/urlinfo/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
