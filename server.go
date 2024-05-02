package main

//will need a HTTP handler
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

import (
	"fmt"
	"net/http"
	"strings"
)

// create struct to assign malware true/false to url
type url struct {
	url     string
	malware bool
}

// will need a HTTP handler
func Server(w http.ResponseWriter, r *http.Request) {
	//trim URL for database lookup
	url := strings.TrimPrefix(r.URL.Path, "/v1/urlinfo/")
	fmt.Fprintf(w, "URL "+url)
	//check database
	//return true/false
}

//lookup database
//func lookup(url string)bool{
//
// }
