package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var reader = strings.NewReader("hi")

var serverTests = []struct {
	name     string
	inMethod string
	inUrl    string
	inReader *strings.Reader
	out      string
}{
	{"Check if URL is not DB", "GET", "/v1/urlinfo/helloworld", nil, "URL helloworld \nURL not found in database"},
	{"Check if URL is in DB", "GET", "/v1/urlinfo/abc.com", nil, "URL abc.com \nMalware present yes"},
	{"Check if not GET method request", "POST", "/v1/urlinfo/abc.com", reader, "Only GET requests are allowed!\n"},
}

func TestGETPage(t *testing.T) {

	for _, tt := range serverTests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.inMethod, tt.inUrl, nil)
			checkHttpErr(err, t)

			response := httptest.NewRecorder()

			Server(response, request)
			got := response.Body.String()
			assertUrlResponse(t, got, tt.out)
		})
	}
}

func assertUrlResponse(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("URL is incorrect - got %q , want %q", got, want)
	}
}

func checkHttpErr(err error, t *testing.T) {
	if err != nil {
		t.Errorf("could not create http request")
	}
}
