package cyoa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var errNoJSONFound = errors.New("No .json files found in this directory")

// StoryArc is used to deseralize a choose-your-own-adventure json
type StoryArc struct {
	Title   string
	Story   []string
	Options []Option
}

// Option represents a user option in a choose-your-own-adventure
type Option struct {
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
// StoryArc structs. The expected JSON format is
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
func parseJSONstory(jsonBytes []byte) (map[string]StoryArc, error) {
	parsedJSON := make(map[string]StoryArc)
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

// addAdventurePages registers handler functions for a choose-your-own-adventure
// story on handler. The handler functions are registered at baseURL/arc-titles.
// baseURL can be, e.g. the story title. The response body is the template
// storyTemplate which has as a top-level context the parsedStory map. See
// the StoryArc struct for details on the exported fields.
func addAdventurePages(handler *http.ServeMux, parsedStory map[string]StoryArc, storyTemplate, baseURL string) error {
	body := template.Must(template.New("StoryArc").Parse(storyTemplate))
	for title, structBody := range parsedStory {
		buf := bytes.NewBuffer(make([]byte, 10))
		err := body.Execute(buf, structBody)
		if err != nil {
			return err
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
	return nil
}

// AdventuresHandler gets all the choose-your-own adventure json files
// in dir. It returns an http mux that, for every parsed json, has
// a handler function at  /filename/arc-title for every arc in the json.
// The html response of the handler function is the storyTemplate.
// storyTemplate HAS TO have a formatting directive %s, which will respectively
// take the filename of each parsed json. When storyTemplate is executed,
// it will have the top-level context be the StoryArc structs corresponding
// to the deserialized jsons from dir.
func AdventuresHandler(dir, storyTemplate string) (*http.ServeMux, error) {
	handler := http.NewServeMux()
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
		path := fmt.Sprintf("/%s", filename)
		templ := fmt.Sprintf(storyTemplate, filename)
		err = addAdventurePages(handler, parsedJSON, templ, path)
		if err != nil {
			return nil, err
		}
	}

	return handler, nil

}

// HomePageHandler returns a handler with a handler function mounted
// at "/". The HTML response is the template templ, which will have as
// a top-level context the storyNames slice.
func HomePageHandler(storyNames []string, templ string) (*http.ServeMux, error) {
	handler := http.NewServeMux()
	homeBody := template.Must(template.New("stories").Funcs(template.FuncMap{"toTitle": toTitle}).Parse(templ))
	buf := bytes.NewBuffer(make([]byte, 10))
	err := homeBody.Execute(buf, storyNames)
	if err != nil {
		return nil, err
	}
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(r.URL.Path))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	})

	return handler, nil
}
