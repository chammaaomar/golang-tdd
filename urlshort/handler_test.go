package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boltdb/bolt"
)

func TestMapHandler(t *testing.T) {
	t.Run("MapHandler should route a request or use fallback if not found", func(t *testing.T) {
		shortenerMap := map[string]string{
			"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
			"/urlshort":       "https://github.com/gophercises/urlshort",
		}
		handler := MapHandler(shortenerMap, defaultHandler())
		testHandler(t, handler, "/urlshort-final", http.StatusFound)
		testHandler(t, handler, "/not/found", http.StatusOK)
	})
}

func TestYAMLHandler(t *testing.T) {
	t.Run("YAMLHandler should route a request or use fallback if not found", func(t *testing.T) {
		yamlDirectory := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
		handler, err := YAMLHandler([]byte(yamlDirectory), defaultHandler())
		if err != nil {
			t.Fatal(err)
		}
		testHandler(t, handler, "/urlshort-final", http.StatusFound)
		testHandler(t, handler, "/not/found", http.StatusOK)
	})
}

func TestJSONHandler(t *testing.T) {
	t.Run("JSONHandler should route a request or use fallback if not found", func(t *testing.T) {
		jsonDirectory := `
[
	{
		"path": "/urlshort",
		"url": "https://github.com/gophercises/urlshort"
	},
	{
		"path": "/urlshort-final",
		"url": "https://github.com/gophercises/urlshort/tree/solution"
	}
]
`
		handler, err := JSONHandler([]byte(jsonDirectory), defaultHandler())
		if err != nil {
			t.Fatal(err)
		}
		testHandler(t, handler, "/urlshort-final", http.StatusFound)
		testHandler(t, handler, "/not/found", http.StatusOK)
	})
}

func TestBoltDBHandler(t *testing.T) {
	t.Run("BoltDBHandler should route a request or use fallback if not found", func(t *testing.T) {
		testDB, err := bolt.Open("test.db", 0600, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer testDB.Close()
		errUpdate := testDB.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucket([]byte("testBucket"))
			if err != nil {
				return err
			}
			err = b.Put([]byte("/urlshort"), []byte("https://github.com/gophercises/urlshort"))
			if err != nil {
				return err
			}
			err = b.Put([]byte("/urlshort-final"), []byte("https://github.com/gophercises/urlshort/tree/solution"))
			if err != nil {
				return err
			}
			return nil
		})
		if errUpdate != nil {
			t.Fatal(errUpdate)
		}
		defer testDB.Update(func(tx *bolt.Tx) error {
			tx.DeleteBucket([]byte("testBucket"))
			return nil
		})
		handler := BoltDBHandler(testDB, "testBucket", defaultHandler())
		testHandler(t, handler, "/urlshort-final", http.StatusFound)
		testHandler(t, handler, "/not/found", http.StatusOK)
	})
}

func testHandler(t *testing.T, handler http.HandlerFunc, path string, expectedStatus int) {
	t.Helper()
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != expectedStatus {
		t.Fatalf("Expected status code %d, got %d", expectedStatus, rr.Code)
	}
}
