package cyoa

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestWebApplication(t *testing.T) {
	t.Run("User should be shown the 'intro' arc when requesting /stories/the-little-blue-gopher", func(t *testing.T) {
		expectedContents := map[string]string{
			"title":     "The Little Blue Gopher",
			"paragraph": "Once upon a time, long long ago, there was a little blue gopher.",
			"option":    "That story about the Sticky Bandits",
		}
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/stories/the-little-blue-gopher/intro", nil)
		if err != nil {
			t.Fatal(err)
		}
		jsonBytes, err := ReadFile("stories/the-little-blue-gopher.json")
		if err != nil {
			t.Fatal(err)
		}
		adventure, err := ParseJSONstory(jsonBytes)
		if err != nil {
			t.Fatal(err)
		}
		storyTempl = fmt.Sprintf(storyTempl, "/stories/the-little-blue-gopher")
		handler, err := AdventureHandler(nil, adventure, storyTempl, "/stories/the-little-blue-gopher")

		if err != nil {
			t.Fatal(err)
		}

		log.Fatal(http.ListenAndServe("localhost:8080", handler))

		handler.ServeHTTP(rr, req)
		if code := rr.Code; code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, code)
		}
		for k, v := range expectedContents {
			if !bytes.Contains(rr.Body.Bytes(), []byte(v)) {
				t.Fail()
				t.Logf("Expected %q to be present, not found", k)
			}
		}
	})
}

func TestParseJSONstory(t *testing.T) {
	jsonBytes := `
{
  "intro": {
    "title": "The Little Blue Gopher",
    "story": [
      "Once upon a time, long long ago, there was a little blue gopher."
    ],
    "options": [
      {
        "text": "That story about the Sticky Bandits isn't real.",
        "arc": "new-york"
      }
    ]
  }
}
`
	parsedStory, err := ParseJSONstory([]byte(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]StoryArc{
		"intro": StoryArc{
			Title: "The Little Blue Gopher",
			Story: []string{
				"Once upon a time, long long ago, there was a little blue gopher.",
			},
			Options: []Option{
				Option{
					Text: "That story about the Sticky Bandits isn't real.",
					Arc:  "new-york",
				},
			},
		},
	}
	if !reflect.DeepEqual(parsedStory, expected) {
		t.Fatalf("Expected parsed story to be %v, got %v", expected, parsedStory)
	}
}
