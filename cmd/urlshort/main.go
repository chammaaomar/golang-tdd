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

func main() {
	flag.Parse()
	mux := defaultMux()
	var routing []byte

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	if len(*yamlPtr) > 0 {
		routing = readFile(*yamlPtr)
	}
	yamlHandler, err := urlshort.YAMLHandler(routing, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on localhost:8080")
	http.ListenAndServe("localhost:8080", yamlHandler)
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
