package cyoa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var errNoJSONFound = errors.New("No .json files found in this directory")

type storyArc struct {
	Title   string
	Story   []string
	Options []option
}

type option struct {
	Text string
	Arc  string
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// parseJSONstory parses a story encoded as a JSON into a Go map of
// storyArc structs. The expected JSON format is
// {
// 		"[ARC TITLE]": {
// 			"title": "[TITLE]",
//			"story": [
//				"[PARAGRAPH]",
//				"[PARAGRAPH]",
//				...
//			],
//			"options": [
//				{
//					"text": "[TEXT]"
//					"arc": "[DYNAMIC ARC TITLE]"
//				}
//			]
//		},
//		...
// }
func parseJSONstory(jsonBytes []byte) (map[string]storyArc, error) {
	parsedJSON := make(map[string]storyArc)
	jsonBuf := bytes.NewBuffer(jsonBytes)
	decoder := json.NewDecoder(jsonBuf)

	for decoder.More() {
		err := decoder.Decode(&parsedJSON)
		if err != nil {
			return nil, err
		}
	}

	return parsedJSON, nil
}

// addAdventurePages takes a handler and adds at the base url handler functions
// for the story arcs provided in the parsed story map. The handler functions
// are registered at baseURL/arc-title, for every story arc.
func addAdventurePages(handler *http.ServeMux, parsedStory map[string]storyArc, storyTemplate, baseURL string) (*http.ServeMux, error) {
	body := template.Must(template.New("StoryArc").Parse(storyTemplate))
	if handler == nil {
		handler = http.NewServeMux()
	}
	for title, structBody := range parsedStory {
		buf := bytes.NewBuffer(make([]byte, 10))
		err := body.Execute(buf, structBody)
		if err != nil {
			log.Fatal(err)
		}
		func(respBody []byte) {
			// copy the template body by value
			handler.HandleFunc(fmt.Sprintf("%s/%s", baseURL, title), func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write(respBody)
			})
		}(buf.Bytes())
		buf.Reset()
	}
	return handler, nil
}

func AddAdventures(handler *http.ServeMux, dir, baseURL, storyTemplate string) (http.Handler, error) {
	if handler == nil {
		handler = http.NewServeMux()
	}
	files, err := filepath.Glob(path.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errNoJSONFound
	}
	for _, file := range files {
		filename := strings.Split(path.Base(file), ".json")[0]
		jsonBytes, err := readFile(file)
		if err != nil {
			return nil, err
		}
		parsedJSON, err := parseJSONstory(jsonBytes)
		if err != nil {
			return nil, err
		}
		path := fmt.Sprintf("%s/%s", baseURL, filename)
		templ := fmt.Sprintf(storyTemplate, path)
		_, err = addAdventurePages(handler, parsedJSON, templ, path)
		if err != nil {
			return nil, err
		}
	}

	return handler, nil

}

// AddStoriesHomePage adds to the handler a handle function for the path
// url. The handle function serves the template templ. The top-level (initial)
//context in the template is the string slice storyNames, which is the only
// restriction on the template.
func AddStoriesHomePage(handler *http.ServeMux, storyNames []string, url, templ string) (*http.ServeMux, error) {
	if handler == nil {
		handler = http.NewServeMux()
	}
	homeBody := template.Must(template.New("stories").Funcs(template.FuncMap{"toTitle": toTitle}).Parse(templ))
	buf := bytes.NewBuffer(make([]byte, 10))
	err := homeBody.Execute(buf, storyNames)
	if err != nil {
		return nil, err
	}
	handler.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	})

	return handler, nil
}
