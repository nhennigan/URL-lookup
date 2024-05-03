package main

//will need a HTTP handler
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

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
	res, err := malwareCheck(url)
	if err != nil {
		fmt.Fprintf(w, "There is some issue with the database")
	}

	if res == "" {
		fmt.Fprintf(w, " \nURL not found in database")
	} else {
		fmt.Fprintf(w, " \nSafe "+res)
	}

	//read entries.json file every 10 mins
	duration := 10 * time.Minute
	ticker := time.NewTicker(duration)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				entries := readNewData()
				addNewEntry(entries)
			}
		}
	}()

	// time.Sleep(5 * time.Minute)
	// ticker.Stop()
	// done <- true
	// for range ticker.C {
	// 	fmt.Printf("ticker happening")
	// 	entries := readNewData()
	// 	addNewEntry(entries)
	// }

}
