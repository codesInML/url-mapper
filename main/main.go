package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	urlshort "github.com/codesInML/url-short"
	"github.com/go-yaml/yaml"
)

func main() {
	yamlFile := flag.String("yaml", "test.yaml", "Load urls from a yaml file. Must be in the accepted yaml format.")
	jsonFile := flag.String("json", "test.json", "Load urls from a json file. Must be in the accepted json format.")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	ymlFile, err := os.Open(*yamlFile)
	if err != nil {
		fmt.Println("Could not load yaml file")
		return
	}

	jsnFile, err := os.Open(*jsonFile)
	if err != nil {
		fmt.Println("Could not load json file")
		return
	}

	json := json.NewDecoder(jsnFile)
	_, err = urlshort.JSONHandler(json, mapHandler)
	if err != nil {
		panic(err)
	}

	yaml := yaml.NewDecoder(ymlFile)
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
