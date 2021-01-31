package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testPathsToUrls = map[string]string{
	"/redirect":         "/netapp",
	"/another-redirect": "/google",
	"/final-redirect":   "/local",
}

var testYaml = `
- path: /urlshort
  url: /my-very-long-url
- path: /urlshort-final
  url: /my-much-much-longer-url
`

var testJson = `
[ 
    {
        "path": "/urlshort",
        "url": "/my-very-long-url"
    },
    {
        "path": "/urlshort-final",
        "url": "/my-much-much-longer-url"
    }
]
`

var testRedirects = []Redirect{
	Redirect{
		Path: "/urlshort",
		URL:  "/my-very-long-url",
	},
	Redirect{
		Path: "/urlshort-final",
		URL:  "/my-much-much-longer-url",
	},
}

// TestMapHandlerRedirect tests given a MapHandler that
// requests for matching routes are correctly redirected
func TestMapHandlerRedirect(t *testing.T) {
	mux := http.NewServeMux()
	handler := MapHandler(testPathsToUrls, mux)

	// Test each mapping to see if a successful redirect occurs
	for path, url := range testPathsToUrls {

		// Create a test response to avoid a real http request
		rr := GetTestResponse(t, handler, path)

		// Validate 302 status coded
		if status := rr.Code; status != http.StatusFound {
			t.Errorf("Expected status code %v, got %v", http.StatusFound, status)
		}

		// Validate redirect location exists and is correct
		if location, ok := rr.HeaderMap["Location"]; ok {
			if url != location[0] {
				t.Errorf("Expected url %s got %s", url, location[0])
			}
		} else {
			t.Errorf("Expected url %s got /", url)
		}
	}
}

// TestMapHandlerFallback tests given a MapHandler that the fallback
// handler is correctly called
func TestMapHandlerFallback(t *testing.T) {
	fallbackPath := "/unused-path"
	fallbackText := "Hello, world"

	mux := http.NewServeMux()
	mux.HandleFunc(fallbackPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fallbackText)
	})

	handler := MapHandler(testPathsToUrls, mux)

	rr := GetTestResponse(t, handler, fallbackPath)

	// Validate 200 status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Validate response body matches our fallback handler text
	if body := rr.Body; body.String() != fallbackText {
		t.Errorf("Expected body %s got %s", fallbackText, body.String())
	}
}

func TestYamlHandlerRedirect(t *testing.T) {
	mux := http.NewServeMux()
	handler, err := YamlHandler([]byte(testYaml), mux)
	if err != nil {
		t.Fatal(err)
	}

	for _, redirect := range testRedirects {
		url := redirect.URL
		path := redirect.Path

		rr := GetTestResponse(t, handler, path)

		if status := rr.Code; status != http.StatusFound {
			t.Errorf("Expected status code %v, got %v", http.StatusFound, status)
		}

		if location, ok := rr.HeaderMap["Location"]; ok {
			if url != location[0] {
				t.Errorf("Expected url %s got %s", url, location[0])
			}
		} else {
			t.Errorf("Expected url %s got /", url)
		}
	}
}

func TestYamlHandlerFallback(t *testing.T) {
	fallbackPath := "/unused-path"
	fallbackText := "Hello, world"

	mux := http.NewServeMux()
	mux.HandleFunc(fallbackPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fallbackText)
	})

	handler, err := YamlHandler([]byte(testYaml), mux)
	if err != nil {
		t.Fatal(err)
	}

	rr := GetTestResponse(t, handler, fallbackPath)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	if body := rr.Body; body.String() != fallbackText {
		t.Errorf("Expected body %s got %s", fallbackText, body.String())
	}
}

func TestJsonHandlerRedirect(t *testing.T) {
	mux := http.NewServeMux()
	handler, err := JsonHandler([]byte(testJson), mux)
	if err != nil {
		t.Fatal(err)
	}

	for _, redirect := range testRedirects {
		url := redirect.URL
		path := redirect.Path

		rr := GetTestResponse(t, handler, path)

		if status := rr.Code; status != http.StatusFound {
			t.Errorf("Expected status code %v, got %v", http.StatusFound, status)
		}

		if location, ok := rr.HeaderMap["Location"]; ok {
			if url != location[0] {
				t.Errorf("Expected url %s got %s", url, location[0])
			}
		} else {
			t.Errorf("Expected url %s got /", url)
		}
	}
}

func TestJsonHandlerFallback(t *testing.T) {
	fallbackPath := "/unused-path"
	fallbackText := "Hello, world"

	mux := http.NewServeMux()
	mux.HandleFunc(fallbackPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fallbackText)
	})

	handler, err := JsonHandler([]byte(testJson), mux)
	if err != nil {
		t.Fatal(err)
	}

	rr := GetTestResponse(t, handler, fallbackPath)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	if body := rr.Body; body.String() != fallbackText {
		t.Errorf("Expected body %s got %s", fallbackText, body.String())
	}
}

func TestParseYamlSuccess(t *testing.T) {
	var redirects []Redirect

	err := ParseYaml([]byte(testYaml), &redirects)
	if err != nil {
		t.Fatalf("Unexpected error from valid yaml")
	}

	if !cmp.Equal(testRedirects, redirects) {
		t.Errorf("Expected redirects %v, got %v", testRedirects, redirects)
	}
}

func TestParseYamlFailure(t *testing.T) {
	invalidYaml := `
- path: /urlshort
	url: https://github.com/gophercises/urlshort
path: /urlshort-final
	url: https://github.com/gophercises/urlshort/tree/solution
`

	var redirects []Redirect

	err := ParseYaml([]byte(invalidYaml), &redirects)
	if err == nil {
		t.Errorf("Expected error from invalid yaml")
	}
}

func TestRedirectsToMap(t *testing.T) {
	pathsToUrls := RedirectsToMap(testRedirects)

	for _, redirect := range testRedirects {
		if url, ok := pathsToUrls[redirect.Path]; ok {
			if url != redirect.URL {
				t.Errorf("Expected path %s got %s", redirect.Path, url)
			}
		} else {
			t.Errorf("Expected to find key %s in map", redirect.Path)
		}
	}
}

func GetTestResponse(t *testing.T, handler http.HandlerFunc, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(rr, req)

	return rr
}
