package urlshort

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/boltdb/bolt"

	"gopkg.in/yaml.v2"
)

type pathURL struct {
	Path string
	URL  string
}

func defaultHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fallbackMessage := ("Hello, you landed at %s. This doesn't correspond to any shortenered URL I know of. Care to register it?")
		fmt.Fprintf(w, fallbackMessage, html.EscapeString(r.URL.Path))
	})
	return mux
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL, ok := pathsToUrls[html.EscapeString(r.URL.Path)]
		if ok {
			http.Redirect(w, r, targetURL, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
		return
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all relate to having
// invalid YAML data.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse yaml into a slice of maps, as appropriate for the expected
	// format
	mappings := make([]map[string]string, 0, 10)
	combinedMapping := make(map[string]string)
	err := yaml.Unmarshal(yml, &mappings)
	if err != nil {
		return nil, err
	}
	// combine the maps in the slice into a single map
	for _, routingPair := range mappings {
		combinedMapping[routingPair["path"]] = routingPair["url"]
	}
	return MapHandler(combinedMapping, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//	[
//		{
//			"path": "/some-path",
//			"url": "https://www.some-url.com/demo"
//		}
//	]
//
// The only errors that can be returned all relate to having
// invalid JSON data.
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	mappings := make([]pathURL, 0, 10)
	jsonReader := bytes.NewReader(jsonData)
	decoder := json.NewDecoder(jsonReader)
	err := decoder.Decode(&mappings)
	if err != nil {
		return nil, err
	}
	mapping := make(map[string]string)
	// combine all mappings extracted from json into single mapping
	for _, routingPair := range mappings {
		mapping[routingPair.Path] = routingPair.URL
	}
	return MapHandler(mapping, fallback), nil
}

// BoltDBHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the BoltDB) to their corresponding URL (values
// that each key in the BoltDB points to, in string format).
// If the path is not provided in the BoltDB, then the fallback
// http.Handler will be called instead.
func BoltDBHandler(db *bolt.DB, bucket string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlBuf := bytes.NewBuffer(make([]byte, 0, 10))
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			if b == nil {
				return bolt.ErrBucketNotFound
			}
			tempURL := b.Get([]byte(html.EscapeString(r.URL.Path)))
			_, err := urlBuf.Write(tempURL)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fallback.ServeHTTP(w, r)
			return
		}
		url := string(urlBuf.Bytes())
		if len(url) == 0 {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
}
