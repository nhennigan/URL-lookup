package main

//will need a HTTP handler
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

// create struct to assign malware true/false to url
// type url struct{
// 	url     string
// 	malware bool
// }

var limiter = rate.NewLimiter(1, 15)

// will need a HTTP handler
func Server(w http.ResponseWriter, r *http.Request) {
	//trim URL for database lookup
	if r.Method != "GET" {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if !limiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	url := strings.TrimPrefix(r.URL.Path, "/v1/urlinfo/")
	fmt.Fprintf(w, "URL "+url)
	//check database
	res := malwareCheck(url)

	if res == "" {
		fmt.Fprintf(w, " URL not found in database")
	} else {
		fmt.Fprintf(w, " Safe "+res)
	}
}
