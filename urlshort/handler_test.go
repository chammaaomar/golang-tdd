package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMapHandler(t *testing.T) {
	t.Run("MapHandler should correctly route a path present in the Map", func(t *testing.T) {
		shortenerMap := map[string]string{
			"/yaml-godoc": "https://godoc.org/gopkg.in/yaml.v2",
		}
		handler := MapHandler(shortenerMap, defaultHandler())
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/yaml-godoc", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusFound {
			t.Fatalf("Expected status code %v, got %d", http.StatusFound, status)
		}
	})

	t.Run("MapHandler should use the fallback handler if the path isn't found", func(t *testing.T) {
		shortenerMap := map[string]string{}
		handler := MapHandler(shortenerMap, defaultHandler())
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/using/fallback", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}
