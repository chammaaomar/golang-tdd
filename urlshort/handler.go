package urlshort

import (
	"fmt"
	"html"
	"net/http"

	"gopkg.in/yaml.v2"
)

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
	mappings := make([]map[string]string, 0, 50)
	combinedMapping := make(map[string]string)
	err := yaml.Unmarshal(yml, &mappings)
	if err != nil {
		return nil, err
	}
	for _, routingPair := range mappings {
		combinedMapping[routingPair["path"]] = routingPair["url"]
	}
	return MapHandler(combinedMapping, fallback), nil
}
