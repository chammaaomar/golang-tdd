package cyoa

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAdventurePages(t *testing.T) {
	t.Run("User should be shown the 'intro' arc when requesting /the-little-blue-gopher/intro", func(t *testing.T) {
		expectedContents := map[string]string{
			"title":     "The Little Blue Gopher",
			"paragraph": "Once upon a time, long long ago, there was a little blue gopher.",
			"option":    "That story about the Sticky Bandits",
		}
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/the-little-blue-gopher/intro", nil)
		if err != nil {
			t.Fatal(err)
		}
		jsonBytes, err := readFile("stories/the-little-blue-gopher.json")
		if err != nil {
			t.Fatal(err)
		}
		adventure, err := parseJSONstory(jsonBytes)
		if err != nil {
			t.Fatal(err)
		}
		testTempl := fmt.Sprintf(storyTempl, "/the-little-blue-gopher")
		handler, err := addAdventurePages(nil, adventure, testTempl, "/the-little-blue-gopher")

		if err != nil {
			t.Fatal(err)
		}

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
	parsedStory, err := parseJSONstory([]byte(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]storyArc{
		"intro": storyArc{
			Title: "The Little Blue Gopher",
			Story: []string{
				"Once upon a time, long long ago, there was a little blue gopher.",
			},
			Options: []option{
				option{
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

func TestHomePage(t *testing.T) {
	t.Run("User should be presented with a list of stories, and upon clicking, should navigate to the correct page", func(t *testing.T) {

		// Make initial request to homepage
		// expect to see list of stories
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		expectedContent := "<li><a href=\"/the-little-blue-gopher/intro\"><p>The Little Blue Gopher</p>"

		handler, err := AddStoriesHomePage(nil, []string{"the-little-blue-gopher"}, "/", fmt.Sprintf(homeTempl, ""))

		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if code := rr.Code; code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, code)
		}

		if !bytes.Contains(rr.Body.Bytes(), []byte(expectedContent)) {
			t.Fatalf("Expected HTML response to contain %v, not found", expectedContent)
		}

		// make subsequent request to a story in the list
		// simulating navigation or user selecting an adventure
		// expect the user to be taken to the page

		_, err = AddAdventures(handler, "stories/", "", storyTempl)

		if err != nil {
			t.Fatal(err)
		}

		req, err = http.NewRequest("GET", "/the-little-blue-gopher/intro", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr.Body.Reset()
		handler.ServeHTTP(rr, req)

		expectedContents := map[string]string{
			"title":     "The Little Blue Gopher",
			"paragraph": "Once upon a time, long long ago, there was a little blue gopher.",
			"option":    "That story about the Sticky Bandits",
		}

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
