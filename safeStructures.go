package main

import "sync"

/*
#######################################
Index map
#######################################
*/

// Index map
type SafeMap struct {
	mu         sync.RWMutex
	items      map[string]int
	indexCount int //number of already used indexes
}

// Thread safe acces of elemtes in map
func (sm *SafeMap) Get(key string) (int, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.items[key]
	return value, ok
}

// Add new word and index to index map
func (sm *SafeMap) AddWord(word string) (idx int) {
	idx = sm.indexCount
	sm.items[word] = idx
	sm.indexCount++
	return
}

/*
#######################################
Frequation matrix
#######################################
*/

// Frequation matrix
type FrequencyMatrix struct {
	mu        sync.Mutex
	matrix    [][]float64 //frequation matrix
	wordCount int         //number of idexed words
	capacity  int         //number of columns in matrix
}

// Thread safe method for double column count of frequation matrix
func (frequencyMatrix *FrequencyMatrix) DoubleMatrixIfSmall() {
	frequencyMatrix.mu.Lock()
	defer frequencyMatrix.mu.Unlock()
	if frequencyMatrix.wordCount < frequencyMatrix.capacity-1 {
		return
	}
	file_count := len(frequencyMatrix.matrix)
	new_matrix := make([][]float64, file_count)
	for i := 0; i < file_count; i++ {
		new_matrix[i] = make([]float64, frequencyMatrix.capacity*2)
		copy(new_matrix[i], frequencyMatrix.matrix[i])
	}
	frequencyMatrix.matrix = new_matrix
	frequencyMatrix.capacity *= 2
}

// Initialization of frequation matrix
func (frequencyMatrix *FrequencyMatrix) Init() {
	for i := 0; i < len(frequencyMatrix.matrix); i++ {
		frequencyMatrix.matrix[i] = make([]float64, frequencyMatrix.capacity)
	}
}

// Get count of nonzero elements in given column
func (frequencyMatrix *FrequencyMatrix) CountOfNonZeroInColumn(columnIdx int) (nonzeroCount int) {
	nonzeroCount = 0
	rowCount := len(frequencyMatrix.matrix)
	for rIdx := 0; rIdx < rowCount; rIdx++ {
		if frequencyMatrix.matrix[rIdx][columnIdx] > 0 {
			nonzeroCount++
		}
	}
	return
}
