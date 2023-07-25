package main

import "sort"

// Struckt representing search results in one folder
type SearchResult struct {
	//Indexes of word in frequation matrix, based on what the search is done
	indexes []int
	//Frequencies of each word in this file
	wordFrequencies []float64
	//How many frequencies are not zero
	nonZeroItemscount int
	//Sum of all frequecies
	wordFreqSum float64
	fileIndex   int
}

// Initialization of search result
func (searchResult *SearchResult) init(indexes []int, fileIndex int) {
	searchResult.indexes = indexes
	searchResult.wordFrequencies = make([]float64, len(searchResult.indexes))
	searchResult.nonZeroItemscount = 0
	searchResult.wordFreqSum = 0
	searchResult.fileIndex = fileIndex
}

// Add frequency to search result on given index
func (searchResult *SearchResult) WriteOnIndex(index int, value float64) {
	previousValue := searchResult.wordFrequencies[index]
	searchResult.wordFreqSum += (value - previousValue)
	if previousValue != 0 {
		searchResult.nonZeroItemscount--
	}
	if value != 0 {
		searchResult.nonZeroItemscount++
	}
	searchResult.wordFrequencies[index] = value
}

// Comparator of two search results
func resultComparator(sr1, sr2 SearchResult) bool {
	if sr1.nonZeroItemscount == sr2.nonZeroItemscount {
		return sr1.wordFreqSum > sr2.wordFreqSum
	}

	return sr1.nonZeroItemscount > sr2.nonZeroItemscount
}

// Get list of files based on relevance of search
func SearchText(words []string, frequencyMatrix *FrequencyMatrix, indexMap *SafeMap, files []string) (filePaths []string) {
	wordIndexes := make([]int, 0)
	for _, w := range words {
		idx, ok := indexMap.Get(w)
		if ok {
			wordIndexes = append(wordIndexes, idx)
		}
	}

	results := make([]SearchResult, len(files))

	for fileIndex := 0; fileIndex < len(files); fileIndex++ {
		results[fileIndex] = frequencyMatrix.getResults(wordIndexes, fileIndex)
	}
	sort.Slice(results, func(i, j int) bool {
		return resultComparator(results[i], results[j])
	})

	filePaths = convertResultsToFileNames(results, files)
	return
}

// Convert ordered file indexes to file names
func convertResultsToFileNames(results []SearchResult, files []string) (fileNames []string) {
	for _, result := range results {
		fileNames = append(fileNames, files[result.fileIndex])
	}
	return
}

// Get tf-idf results of given words from frequation matrix
func (frequencyMatrix *FrequencyMatrix) getResults(wordIndexes []int, fileIndex int) (searchResult SearchResult) {
	searchResult.init(wordIndexes, fileIndex)
	for idx, wordIdx := range searchResult.indexes {
		searchResult.WriteOnIndex(idx, frequencyMatrix.matrix[fileIndex][wordIdx])
	}
	return
}
