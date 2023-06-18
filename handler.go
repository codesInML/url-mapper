package urlshort

import (
	"encoding/json"
	"net/http"

	"github.com/go-yaml/yaml"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml *yaml.Decoder, fallback http.Handler) (http.HandlerFunc, error) {
	urlPaths, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMap(urlPaths)

	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(yml *json.Decoder, fallback http.Handler) (http.HandlerFunc, error) {
	urlPaths, err := parseJson(yml)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMap(urlPaths)

	return MapHandler(pathsToUrls, fallback), nil
}

func parseJson(data *json.Decoder) ([]urlPath, error) {
	var urlPaths []urlPath
	err := data.Decode(&urlPaths)

	if err != nil {
		return nil, err
	}

	return urlPaths, nil
}

func parseYaml(data *yaml.Decoder) ([]urlPath, error) {
	var urlPaths []urlPath
	err := data.Decode(&urlPaths)

	if err != nil {
		return nil, err
	}

	return urlPaths, nil
}

func buildMap(urlPaths []urlPath) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, path := range urlPaths {
		pathsToUrls[path.Path] = path.Url
	}

	return pathsToUrls
}

type urlPath struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
