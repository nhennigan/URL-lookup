package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 15)

// HTTP handler
func Server(w http.ResponseWriter, r *http.Request) {
	// if HTTP method is not GET, do not allow the request
	if r.Method != "GET" {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	//only allow a limited number of requests per second
	if !limiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	//trim URL for database lookup
	url := strings.TrimPrefix(r.URL.Path, "/v1/urlinfo/")
	fmt.Fprintf(w, "URL "+url)
	//check database for malware
	res, err := malwareCheck(url)
	if err != nil {
		fmt.Fprintf(w, "There is some issue with the database")
	}

	//print result to screen
	if res == "" {
		fmt.Fprintf(w, " \nURL not found in database")
	} else {
		fmt.Fprintf(w, " \nMalware present "+res)
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
}
