package urlshort

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type Redirect struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
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
// The only errors that can be returned all related to having
// invalid YAML data.
func YamlHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirects []Redirect

	err := ParseYaml(yml, &redirects)
	if err != nil {
		return nil, err
	}

	pathsToUrls := RedirectsToMap(redirects)

	return MapHandler(pathsToUrls, fallback), nil
}

// JsonHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
// [
//     {
//         "path": "/urlshort",
//         "url": "https://github.com/gophercises/urlshort"
//     },
//     {
//         "path": "/urlshort-final",
//         "url": "https://github.com/gophercises/urlshort/tree/solution"
//     }
// ]
// The only errors that can be returned all related to having
// invalid json data.
func JsonHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirects []Redirect

	err := ParseJson(json, &redirects)
	if err != nil {
		return nil, err
	}

	pathsToUrls := RedirectsToMap(redirects)

	return MapHandler(pathsToUrls, fallback), nil
}

// DbHandler retrieves path -> url mappings from a provided
// database and then returns an http.HandlerFunc that will
// attempt to map any paths to their corresponding URL. If
// the path is not provided in the JSON, then thefallback
// http.Handler will be called instead.
//
// Expects a database result of []Redirect
// Errors returned are related to the query of the database
// or the unmarshalling of the result into our []Redirect
func DbHandler(ctx context.Context, r *RedirectoryDatabase, fallback http.Handler) (http.HandlerFunc, error) {
	var redirects []Redirect

	// Find all redirect collection items using bson.M{}
	err := r.Find(ctx, &redirects, bson.M{})
	if err != nil {
		return nil, err
	}

	pathsToUrls := RedirectsToMap(redirects)

	return MapHandler(pathsToUrls, fallback), nil
}

// ParseYaml takes raw yaml bytes array and parses
// it into a given interface
func ParseYaml(yml []byte, v interface{}) error {
	return yaml.Unmarshal(yml, v)
}

// ParseJson takes raw json bytes array and parses
// it into a given interface
func ParseJson(jsn []byte, v interface{}) error {
	return json.Unmarshal(jsn, v)
}

// ParseYamlFile takes a path to a yaml file and
// parses it into a given interface
func ParseYamlFile(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	reader := bufio.NewReader(f)

	yamlDecoder := yaml.NewDecoder(reader)

	err = yamlDecoder.Decode(v)

	return err
}

// RedirectsToMap creates a URL[Path] map from an array of
// type Redirect
func RedirectsToMap(redirects []Redirect) map[string]string {
	pathsToUrls := map[string]string{}
	for _, redirect := range redirects {
		pathsToUrls[redirect.Path] = redirect.URL
	}

	return pathsToUrls
}

// RedirectsToArrayInterface creates an array of type interface
// from an array of type Redirect
func RedirectsToArrayInterface(redirects []Redirect) []interface{} {
	s := make([]interface{}, len(redirects))
	for i, v := range redirects {
		s[i] = v
	}

	return s
}
