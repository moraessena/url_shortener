package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/danielAang/url_shortener/urlshort"
)

func main() {
	yamlPath := flag.String("yml", "./files/urls.yml", "Path to the yaml file")
	jsonPath := flag.String("json", "./files/urls.json", "Path to the json file")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World!")
	})
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	yamlHandler, err := urlshort.YAMLHandler(*yamlPath, mapHandler)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshort.JsonHandler(*jsonPath, yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}
