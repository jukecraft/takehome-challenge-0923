package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

const fileDirectory = "./static"
const basePath = "/"
const urlForSearchFunction = "/search"
const environmentVariableForPort = "PORT"
const empty = ""
const defaultPort = "3001"
const queryUrlParameter = "q"
const contentTypeHeader = "Content-Type"
const contentTypeJson = "application/json"
const resultWindow = 250
const maxResults = 20
const filenameToSearchIn = "completeworks.txt"
const regexForCaseInsensitiveSearch = "(?i)"
const logMessageForSearchAvailable = "shakesearch available at http://localhost:%s..."
const logMessageForPort = ":%s"
const errorMessageSearchQueryMissing = "missing search query in URL params"
const errorMessageEncodingFailure = "encoding failure"
const errorMessageWritingFailure = "error writing: %v"
const errorMessageForLoadFailure = "load: %w"

func main() {
	searcher := loadCompleteWorksOfShakespeare()
	setUpFileServer()
	setUpSearchHandler(searcher)
}

func loadCompleteWorksOfShakespeare() Searcher {
	searcher := Searcher{}
	err := searcher.Load(filenameToSearchIn)
	if err != nil {
		log.Fatal(err)
	}
	return searcher
}

func setUpSearchHandler(searcher Searcher) {
	http.HandleFunc(urlForSearchFunction, handleSearch(searcher))

	port := os.Getenv(environmentVariableForPort)
	if port == empty {
		port = defaultPort
	}

	fmt.Printf(logMessageForSearchAvailable, port)
	err := http.ListenAndServe(fmt.Sprintf(logMessageForPort, port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setUpFileServer() {
	fs := http.FileServer(http.Dir(fileDirectory))
	http.Handle(basePath, fs)
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()[queryUrlParameter]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			writeWithErrorHandling(w, []byte(errorMessageSearchQueryMissing))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeWithErrorHandling(w, []byte(errorMessageEncodingFailure))
			return
		}
		w.Header().Set(contentTypeHeader, contentTypeJson)
		writeWithErrorHandling(w, buf.Bytes())
	}
}

func writeWithErrorHandling(w http.ResponseWriter, bytesToWrite []byte) {
	_, err := w.Write(bytesToWrite)
	if err != nil {

		log.Printf(errorMessageWritingFailure, err)
	}
}

func (searcher *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(errorMessageForLoadFailure, err)
	}
	searcher.CompleteWorks = string(dat)
	searcher.SuffixArray = suffixarray.New(dat)
	return nil
}

func (searcher *Searcher) Search(query string) []string {
	caseInsensitiveSearch := regexp.MustCompile(regexForCaseInsensitiveSearch + query)
	indexesOfFoundOccurrences := searcher.SuffixArray.FindAllIndex(caseInsensitiveSearch, maxResults)
	var results []string
	for _, startAndEndIndex := range indexesOfFoundOccurrences {
		startIndex := startAndEndIndex[0]
		results = append(results, searcher.CompleteWorks[startIndex-resultWindow:startIndex+resultWindow])
	}
	return results
}
