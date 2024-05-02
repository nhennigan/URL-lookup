package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPage(t *testing.T) {
	t.Run("returns check url", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/v1/ulinfo", nil)
		if err != nil {
			t.Errorf("could not create http request")
		}
		response := httptest.NewRecorder()

		handler(response, request)
		got := response.Body.String()
		want := "Check URL"

		if got != want {
			t.Errorf("got %q , want %q", got, want)
		}
	})

}
