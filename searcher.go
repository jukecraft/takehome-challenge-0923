package main

import (
	"index/suffixarray"
	"regexp"
)

const resultWindow = 250
const maxNewResults = 20
const regexForCaseInsensitiveSearch = "(?i)"

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
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
