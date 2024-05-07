package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/json"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 15)

type URL struct {
	URL     string `json:"URL"`
	Malware string `json:"Malware"`
}

type auth struct {
	uname string
	pw    string
}

// HTTP handler
func Server(w http.ResponseWriter, r *http.Request) {
	//check password in request header is correct
	check := auth{os.Getenv("appUsername"), os.Getenv("appPassword")}
	username, password, ok := r.BasicAuth()
	if ok {
		if username == check.uname && password == check.pw {
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

			w.Header().Set("Content-Type", "application/json")
			//trim URL for database lookup
			url := strings.TrimPrefix(r.URL.Path, "/v1/urlinfo/")
			//check database for malware
			res, err := malwareCheck(url)
			if err != nil {
				fmt.Fprintf(w, "There is some issue with the database - %s", err)
			}

			var urlResp URL
			urlResp.URL = url
			//return json response
			if res == "" {
				urlResp.Malware = "unknown"
			} else {
				urlResp.Malware = res
			}

			json.NewEncoder(w).Encode(urlResp)
		}

	} else {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
				entries, err := readNewData()
				if err == nil {
					addNewEntry(entries)
				}

			}
		}
	}()
}
