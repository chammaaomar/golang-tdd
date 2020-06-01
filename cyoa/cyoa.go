package cyoa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var storyTempl = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Choose Your Own Adventure!</title>
	</head>
	<body>
		
		<h1>{{.Title}}</h1>
		{{range .Story}}
		<p>{{.}}</p>
		{{end}}

		<ul>
			{{range .Options}}
			<li>
				<a href="%s/{{.Arc}}"><p>{{.Text}}</p>
			</li>
			{{end}}
		</ul>
		</body>
</html>
`

type StoryArc struct {
	Title   string
	Story   []string
	Options []Option
}

type Option struct {
	Text string
	Arc  string
}

func ReadFile(path string) ([]byte, error) {
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

// ParseJSONstory parses a story encoded as a JSON into a Go map of
// StoryArc structs. The expected JSON format is
// {
// 		"[DYNAMIC ARC TITLE]": {
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
func ParseJSONstory(jsonBytes []byte) (map[string]StoryArc, error) {
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

// AdventureHandler does stuff
func AdventureHandler(handler *http.ServeMux, parsedStory map[string]StoryArc, storyTemplate, baseURL string) (*http.ServeMux, error) {
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
