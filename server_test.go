package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGETPage(t *testing.T) {
	t.Run("check only GET method allowed", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, "/v1/urlinfo/helloworld", strings.NewReader("hi"))
		if err != nil {
			t.Errorf("could not create http request")
		}
		response := httptest.NewRecorder()

		Server(response, request)
		got := response.Body.String()
		want := "Only GET requests are allowed!\n"

		assertUrlResponse(t, got, want)
	})

	t.Run("check URL not found", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/urlinfo/helloworld", nil)
		if err != nil {
			t.Errorf("could not create http request")
		}
		response := httptest.NewRecorder()

		Server(response, request)
		got := response.Body.String()
		want := "URL helloworld \nURL not found in database"

		assertUrlResponse(t, got, want)
	})
	// t.Run("returns check url helloworld", func(t *testing.T) {
	// 	request, err := http.NewRequest(http.MethodGet, "/v1/urlinfo/helloworld", nil)
	// 	if err != nil {
	// 		t.Errorf("could not create http request")
	// 	}
	// 	response := httptest.NewRecorder()

	// 	Server(response, request)
	// 	got := response.Body.String()
	// 	want := "URL helloworld Safe "

	// 	assertUrlResponse(t, got, want)
	// })

	// t.Run("returns check url abcd", func(t *testing.T) {
	// 	request, err := http.NewRequest(http.MethodGet, "/v1/urlinfo/abcd", nil)
	// 	if err != nil {
	// 		t.Errorf("could not create http request")
	// 	}
	// 	response := httptest.NewRecorder()

	// 	Server(response, request)
	// 	got := response.Body.String()
	// 	want := "URL abcd Safe "

	// 	assertUrlResponse(t, got, want)
	// })

	t.Run("returns check url abc.com and safe or not", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/urlinfo/abc.com", nil)
		if err != nil {
			t.Errorf("could not create http request")
		}
		response := httptest.NewRecorder()

		Server(response, request)
		got := response.Body.String()
		want := "URL abc.com \nSafe yes"

		assertUrlResponse(t, got, want)
	})

}

func assertUrlResponse(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("URL is incorrect - got %q , want %q", got, want)
	}
}
