package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIRouteImage(t *testing.T) {
	basePath := ".."
	r := router(basePath, nil)
	ts := httptest.NewServer(r)
	defer ts.Close()

	// issue can't find templates/*.html
	t.Run("GET /image/maze", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/image/maze")
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
		}
		defer resp.Body.Close()
	})
}
