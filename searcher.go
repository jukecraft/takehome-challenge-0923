package main

import (
	"bytes"
	"encoding/json"
	"index/suffixarray"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

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

func (searcher *Searcher) Search(query string, existing int) []string {
	caseInsensitiveSearch := regexp.MustCompile(regexForCaseInsensitiveSearch + query)
	indexesOfFoundOccurrences := searcher.SuffixArray.FindAllIndex(caseInsensitiveSearch, -1)
	endIndex := min(len(indexesOfFoundOccurrences), existing+maxNewResults)
	return collectResults(indexesOfFoundOccurrences[:endIndex], searcher)
}
