package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchHamlet(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "Hamlet"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearchRequest(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, result := range results {
		if strings.Contains(strings.ToLower(result), strings.ToLower(query)) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected result not found for query: %s", query)
	}
}

func TestSearchCaseSensitive(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "hAmLeT"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearchRequest(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, result := range results {
		if strings.Contains(strings.ToLower(result), strings.ToLower(query)) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected result not found for query: %s", query)
	}
}

func TestSearchDrunk(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "drunk"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearchRequest(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 20 {
		t.Errorf("expected 20 results for query: %s, got %d", query, len(results))
	}
}

func TestSearchWithMoreResults(test *testing.T) {
	searcher := Searcher{}
	potentialError := searcher.Load("completeworks.txt")
	if potentialError != nil {
		test.Fatal(potentialError)
	}

	results := getResultsFromQuery(test, potentialError, "/search?q=drunk", searcher)
	moreResults := getResultsFromQuery(test, potentialError, "/search?q=drunk&existing=20", searcher)

	if !(len(moreResults) > len(results)) {
		test.Errorf("expected larger number of results for loading more results for drunk")
	}
}

func getResultsFromQuery(test *testing.T, potentialError error, queryCall string, searcher Searcher) []string {
	request, potentialError := http.NewRequest("GET", queryCall, nil)
	if potentialError != nil {
		test.Fatal(potentialError)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearchRequest(searcher))
	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		test.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	potentialError = json.Unmarshal(responseRecorder.Body.Bytes(), &results)
	if potentialError != nil {
		test.Fatal(potentialError)
	}
	return results
}
