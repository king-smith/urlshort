# Exercise #2: URL Shortener

## Usage
Install deps 

`go dep`

Run 

`go run main/main.go --task=N --path=P`
  - `task`: Task # to run (1 - 5)
  - `path`: Path to file to parse (Tasks 2, 3, 5)
  
 ### Task 4
 Task 4 requires a running mongoDB server to successfully run. This can be done with `docker`:
 
 `docker run -d -p 2701:27017 --name mongodb mongo:4.0.4` 
 
 Create .env file with default mongoDB config

`cat .env.example > .env`


## Exercise details

The goal of this exercise is to create an [http.Handler](https://golang.org/pkg/net/http/#Handler) that will look at the path of any incoming web request and determine if it should redirect the user to a new page, much like URL shortener would.

For instance, if we have a redirect setup for `/dogs` to `https://www.somesite.com/a-story-about-dogs` we would look for any incoming web requests with the path `/dogs` and redirect them.

To complete this exercises you will need to implement the stubbed out methods in [handler.go](https://github.com/gophercises/urlshort/blob/master/handler.go). There are a good bit of comments explaining what each method should do, and there is also a [main/main.go](https://github.com/gophercises/urlshort/blob/master/main/main.go) source file that uses the package to help you test your code and get an idea of what your program should be doing.

I suggest first commenting out all of the code in main.go related to the `YamlHandler` function and focusing on implementing the `MapHandler` function first.

Once you have that working, focus on parsing the YAML using the [gopkg.in/yaml.v2](https://godoc.org/gopkg.in/yaml.v2) package. *Note: You will need to `go get` this package if you don't have it already.*

After you get the YAML parsing down, try to convert the data into a map and then use the MapHandler to finish the YamlHandler implementation. Eg you might end up with some code like this:

```go
func YamlHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
  parsedYaml, err := parseYAML(yaml)
  if err != nil {
    return nil, err
  }
  pathMap := buildMap(parsedYaml)
  return MapHandler(pathMap, fallback), nil
}
```

But in order for this to work you will need to create functions like `parseYAML` and `buildMap` on your own. This should give you ample experience working with YAML data.


## Bonus

As a bonus exercises you can also...

1. Update the [main/main.go](https://github.com/gophercises/urlshort/blob/master/main/main.go) source file to accept a YAML file as a flag and then load the YAML from a file rather than from a string.
2. Build a JSONHandler that serves the same purpose, but reads from JSON data.
3. Build a Handler that doesn't read from a map but instead reads from a database. Whether you use BoltDB, SQL, or something else is entirely up to you.
