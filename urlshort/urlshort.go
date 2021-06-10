package urlshort

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if url, ok := pathsToUrls[p]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func InMemoryHandler(fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func YAMLHandler(filePath string, fallback http.Handler) (http.HandlerFunc, error) {
	var urls []pathUrl
	file := fileFromPath(filePath)
	if err := yaml.Unmarshal(file, &urls); err != nil {
		return nil, err
	}
	paths := buildUrlMap(urls)
	return MapHandler(paths, fallback), nil
}

func JsonHandler(filePath string, fallback http.Handler) (http.HandlerFunc, error) {
	var urls []pathUrl
	file := fileFromPath(filePath)
	if err := json.Unmarshal(file, &urls); err != nil {
		return nil, err
	}
	paths := buildUrlMap(urls)
	return MapHandler(paths, fallback), nil
}

func fileFromPath(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return file
}

func buildUrlMap(urls []pathUrl) map[string]string {
	paths := make(map[string]string, len(urls))
	for _, value := range urls {
		paths[value.Path] = value.URL
	}
	return paths
}
