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
	"strconv"
)

const fileDirectory = "./static"
const basePath = "/"
const urlForSearchFunction = "/search"
const environmentVariableForPort = "PORT"
const empty = ""
const defaultPort = "3001"
const queryUrlParameter = "q"
const existingUrlParameter = "existing"
const contentTypeHeader = "Content-Type"
const contentTypeJson = "application/json"
const filenameToSearchIn = "completeworks.txt"

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
	http.HandleFunc(urlForSearchFunction, handleSearchRequest(searcher))

	port := os.Getenv(environmentVariableForPort)
	if port == empty {
		port = defaultPort
	}

	fmt.Printf("shakesearch available at http://localhost:%s...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setUpFileServer() {
	fileServer := http.FileServer(http.Dir(fileDirectory))
	http.Handle(basePath, fileServer)
}

func (searcher *Searcher) Load(filename string) error {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	searcher.CompleteWorks = string(fileContent)
	searcher.SuffixArray = suffixarray.New(fileContent)
	return nil
}

func handleSearchRequest(searcher Searcher) func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		query, ok := request.URL.Query()[queryUrlParameter]
		if !ok || len(query[0]) < 1 {
			responseWriter.WriteHeader(http.StatusBadRequest)
			write(responseWriter, []byte(("missing search query in URL params")))
			return
		}
		var existing = 0

		existingFromUrl := request.URL.Query()[existingUrlParameter]
		if len(existingFromUrl) != 0 {
			potentialExisting, err := strconv.Atoi(existingFromUrl[0])
			if err != nil {
				responseWriter.WriteHeader(http.StatusBadRequest)
				write(responseWriter, []byte(("expecting existing to be parseable into an integer")))
				return
			}
			existing = potentialExisting
		}

		results := searcher.Search(query[0], existing)
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		err := encoder.Encode(results)
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			write(responseWriter, []byte(("encoding failure")))
			return
		}
		responseWriter.Header().Set(contentTypeHeader, contentTypeJson)
		write(responseWriter, buffer.Bytes())
	}
}

func write(writer http.ResponseWriter, bytesToWrite []byte) {
	_, err := writer.Write(bytesToWrite)
	if err != nil {
		log.Printf("error writing: %v", err)
	}
}
