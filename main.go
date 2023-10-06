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
	potentialError := searcher.Load(filenameToSearchIn)
	if potentialError != nil {
		log.Fatal(potentialError)
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
	potentialError := http.ListenAndServe(fmt.Sprintf(logMessageForPort, port), nil)
	if potentialError != nil {
		log.Fatal(potentialError)
	}
}

func setUpFileServer() {
	fileServer := http.FileServer(http.Dir(fileDirectory))
	http.Handle(basePath, fileServer)
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearchRequest(searcher Searcher) func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		query, ok := request.URL.Query()[queryUrlParameter]
		if !ok || len(query[0]) < 1 {
			responseWriter.WriteHeader(http.StatusBadRequest)
			write(responseWriter, []byte(errorMessageSearchQueryMissing))
			return
		}
		var existing = 0

		existingFromUrl := request.URL.Query()[existingUrlParameter]
		if len(existingFromUrl) != 0 {
			potentialExisting, potentialError := strconv.Atoi(existingFromUrl[0])
			if potentialError != nil {
				responseWriter.WriteHeader(http.StatusBadRequest)
				write(responseWriter, []byte(errorMessageExistingMalformed))
				return
			}
			existing = potentialExisting
		}

		results := searcher.Search(query[0], existing)
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		potentialError := encoder.Encode(results)
		if potentialError != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			write(responseWriter, []byte(errorMessageEncodingFailure))
			return
		}
		responseWriter.Header().Set(contentTypeHeader, contentTypeJson)
		write(responseWriter, buffer.Bytes())
	}
}

func write(writer http.ResponseWriter, bytesToWrite []byte) {
	_, potentialError := writer.Write(bytesToWrite)
	if potentialError != nil {
		log.Printf(errorMessageWritingFailure, potentialError)
	}
}

func (searcher *Searcher) Load(filename string) error {
	fileContent, potentialError := ioutil.ReadFile(filename)
	if potentialError != nil {
		return fmt.Errorf(errorMessageForLoadFailure, potentialError)
	}
	searcher.CompleteWorks = string(fileContent)
	searcher.SuffixArray = suffixarray.New(fileContent)
	return nil
}

func (searcher *Searcher) Search(query string, existing int) []string {
	caseInsensitiveSearch := regexp.MustCompile(regexForCaseInsensitiveSearch + query)
	indexesOfFoundOccurrences := searcher.SuffixArray.FindAllIndex(caseInsensitiveSearch, -1)
	endIndex := min(len(indexesOfFoundOccurrences), existing+maxNewResults)
	return collectResults(indexesOfFoundOccurrences[:endIndex], searcher)
}
func min(first, second int) int {
	if first < second {
		return first
	}
	return second
}
func collectResults(indexesToReturn [][]int, searcher *Searcher) []string {
	var results []string
	for _, startAndEndIndex := range indexesToReturn {
		startIndex := startAndEndIndex[0]
		results = append(results, searcher.CompleteWorks[startIndex-resultWindow:startIndex+resultWindow])
	}
	return results
}
