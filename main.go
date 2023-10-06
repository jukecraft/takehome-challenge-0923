package main

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
const resultWindow = 250
const maxNewResults = 20
const filenameToSearchIn = "completeworks.txt"
const regexForCaseInsensitiveSearch = "(?i)"
const logMessageForSearchAvailable = "shakesearch available at http://localhost:%s..."
const logMessageForPort = ":%s"
const errorMessageSearchQueryMissing = "missing search query in URL params"
const errorMessageExistingMalformed = "expecting existing to be parseable into an integer"
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
	http.HandleFunc(urlForSearchFunction, handleSearchRequest(searcher))

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
	fileServer := http.FileServer(http.Dir(fileDirectory))
	http.Handle(basePath, fileServer)
}

func (searcher *Searcher) Load(filename string) error {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(errorMessageForLoadFailure, err)
	}
	searcher.CompleteWorks = string(fileContent)
	searcher.SuffixArray = suffixarray.New(fileContent)
	return nil
}
