package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chammaaomar/golang-tdd/urlshort"
)

var yamlPtr = flag.String("yaml", "", "path to yaml file mapping short paths to full urls")
var jsonPtr = flag.String("json", "", "path to json file mapping short paths to full urls. Ignored if yaml flag is used.")

func main() {
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	if len(*yamlPtr) > 0 {
		startHandler(*yamlPtr, "yaml", mapHandler)
	} else if len(*jsonPtr) > 0 {
		startHandler(*jsonPtr, "json", mapHandler)
	} else {
		log.Fatal("One of yaml or json flags is required")
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})
	return mux
}

func readFile(path string) []byte {
	file, errOpen := os.Open(path)
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	routings, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return routings
}

// startHandler invokes the relevant handler based on the routing file
// the user chose, i.e. json or yaml.
func startHandler(path string, dataType string, fallback http.HandlerFunc) {
	routing := readFile(path)
	var handler http.HandlerFunc
	var err error
	switch dataType {
	case "yaml":
		handler, err = urlshort.YAMLHandler(routing, fallback)
	case "json":
		handler, err = urlshort.JSONHandler(routing, fallback)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on localhost:8080")
	http.ListenAndServe("localhost:8080", handler)
}
