package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPage(t *testing.T) {
	t.Run("returns check url helloworld", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/urlinfo/helloworld", nil)
		if err != nil {
			t.Errorf("could not create http request")
		}
		response := httptest.NewRecorder()

		Server(response, request)
		got := response.Body.String()
		want := "URL helloworld"

		assertUrlResponse(t, got, want)
	})

	t.Run("returns check url abcd", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/urlinfo/abcd", nil)
		if err != nil {
			t.Errorf("could not create http request")
		}
		response := httptest.NewRecorder()

		Server(response, request)
		got := response.Body.String()
		want := "URL abcd"

		assertUrlResponse(t, got, want)
	})

}

func assertUrlResponse(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("URL is incorrect - got %q , want %q", got, want)
	}
}
