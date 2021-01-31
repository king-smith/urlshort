package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/king-smith/urlshort"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	pathPtr := flag.String("path", "", "Path to input file")
	taskPtr := flag.Int("task", 1, "Task to run")

	flag.Parse()

	switch task := *taskPtr; task {
	case 1:
		Task1()
	case 2:
		Task2(*pathPtr)
	case 3:
		Task3(*pathPtr)
	case 4:
		Task4()
	default:
		Task1()
	}

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func Task1() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YamlHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YamlHandler([]byte(yaml), mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func Task2(path string) {
	// Read full file from provided path
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	mux := defaultMux()

	// Create yaml handler from raw []byte yaml with mux fallback
	yamlHandler, err := urlshort.YamlHandler(yml, mux)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func Task3(path string) {
	// Read full file from provided path
	jsn, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	mux := defaultMux()

	// Create json handler from raw []byte json with mux fallback
	jsonHandler, err := urlshort.JsonHandler(jsn, mux)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func Task4() {
	// Load in .env variables for secret variable safety
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create mongoDB client
	mongoURI := fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Create timeout context
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// Connect to mongoDB client
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Defer disconnect until function ends
	defer client.Disconnect(ctx)

	// Ping mongoDB server to test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)

	// Reset collection for the sake of this exercise
	db.Collection("redirect").Drop(ctx)

	// Create our db struct
	r := urlshort.NewRedirectoryDatabase(db)

	redirects := []urlshort.Redirect{
		urlshort.Redirect{
			Path: "/urlshort",
			URL:  "https://github.com/gophercises/urlshort",
		},
		urlshort.Redirect{
			Path: "/urlshort-final",
			URL:  "https://github.com/gophercises/urlshort/tree/solution",
		},
	}

	// Create array interface from our []Redirect object
	items := urlshort.RedirectsToArrayInterface(redirects)

	// Insert our redirections into the collection
	err = r.InsertMany(ctx, items)
	if err != nil {
		log.Fatal(err)
	}

	mux := defaultMux()

	// Create handler which looks up our redirections in the collection
	dbHandler, err := urlshort.DbHandler(ctx, r, mux)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}
