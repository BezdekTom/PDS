package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// Add record about word to index map and frequation matrix
func indexWord(word string, fileIndex int, idxMap *SafeMap, freqMatrix *FrequencyMatrix) {
	idx, ok := idxMap.Get(word)
	if ok {
		freqMatrix.matrix[fileIndex][idx]++
		return
	}

	idxMap.mu.Lock()
	defer idxMap.mu.Unlock()
	idx, ok = idxMap.items[word]
	if ok {
		freqMatrix.matrix[fileIndex][idx]++
		return
	}
	idx = idxMap.AddWord(word)

	if freqMatrix.capacity <= freqMatrix.wordCount {
		freqMatrix.DoubleMatrixIfSmall()
	}
	freqMatrix.matrix[fileIndex][idx]++
	freqMatrix.wordCount++
}

// Get TF result of given file
func computeTf(filename string, fileIndex int, indexMap *SafeMap, freqMatrix *FrequencyMatrix) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	textsToIgnore := []string{",", ".", ";", ")", "(", "]", "[", "\"", "!", "?"}

	wordCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := SplitTextOnWords(line, textsToIgnore)
		for _, w := range words {
			indexWord(w, fileIndex, indexMap, freqMatrix)
			wordCount++
		}

	}

	convertCountToFrequency(fileIndex, freqMatrix, wordCount)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Compute TF of multiple files
func computeMultipleTf(indexMap *SafeMap, freqMatrix *FrequencyMatrix, files []string, folder string, startIndex int, endIndex int, wg *sync.WaitGroup) {
	for fileIndex := startIndex; fileIndex < endIndex; fileIndex++ {
		filePath := path.Join(folder, files[fileIndex])
		computeTf(filePath, fileIndex, indexMap, freqMatrix)
	}
	wg.Done()
}

// Convert count of word in one file to word frequency
func convertCountToFrequency(fileIndex int, frequencyMatrix *FrequencyMatrix, wordCount int) {
	length := len(frequencyMatrix.matrix[fileIndex])
	for i := 0; i < length; i++ {
		frequencyMatrix.matrix[fileIndex][i] /= float64(wordCount)
	}
}

// Compute IDF of all words
func computeIdf(frequencyMatrix *FrequencyMatrix, idfs []float64, wordStartIdx int, wordEndIdx int, wg *sync.WaitGroup) {
	fileCount := len(frequencyMatrix.matrix)

	for wIdx := wordStartIdx; wIdx < wordEndIdx; wIdx++ {
		idfs[wIdx] = math.Log10(float64(frequencyMatrix.CountOfNonZeroInColumn(wIdx))/float64(fileCount)) + 1
	}
	wg.Done()
}

// Put together information about TF and IDF to TF-IDF
func combineTfAndIidf(frequencyMatrix *FrequencyMatrix, idfs []float64, fileStartIdx int, fileEndIdx int, wg *sync.WaitGroup) {
	for fIdx := fileStartIdx; fIdx < fileEndIdx; fIdx++ {
		for wIdx := 0; wIdx < frequencyMatrix.wordCount; wIdx++ {
			frequencyMatrix.matrix[fIdx][wIdx] *= idfs[fIdx]
		}
	}
	wg.Done()
}

// Compute TF-IDF over all files
func ComputeTfIdf(folder string, threadsCount int) (indexMap *SafeMap, freqMatrix *FrequencyMatrix, files []string) {
	files = ListDir(folder)

	indexMap = &SafeMap{
		indexCount: 0,
		items:      map[string]int{},
	}

	freqMatrix = &FrequencyMatrix{
		wordCount: 0,
		capacity:  4,
		matrix:    make([][]float64, len(files)),
	}
	freqMatrix.Init()

	threadsCount = int(math.Min(float64(threadsCount), float64(len(files))))

	fileIntervalLength := int(len(files) / threadsCount)

	start := time.Now()

	//Compute tf for all documents
	//Certain amount of files on one thread
	var wg sync.WaitGroup
	for i := 0; i < threadsCount-1; i++ {
		startIndex := i * fileIntervalLength
		endIndex := (i + 1) * fileIntervalLength
		wg.Add(1)
		go computeMultipleTf(indexMap, freqMatrix, files, folder, startIndex, endIndex, &wg)
	}
	wg.Add(1)
	go computeMultipleTf(indexMap, freqMatrix, files, folder, (threadsCount-1)*fileIntervalLength, len(files), &wg)
	wg.Wait()

	//Compute idf
	wordIntervalLength := int(freqMatrix.wordCount / threadsCount)
	idfs := make([]float64, freqMatrix.wordCount)
	for i := 0; i < threadsCount-1; i++ {
		startIndex := i * wordIntervalLength
		endIndex := (i + 1) * wordIntervalLength
		wg.Add(1)
		go computeIdf(freqMatrix, idfs, startIndex, endIndex, &wg)
	}
	wg.Add(1)
	go computeIdf(freqMatrix, idfs, (threadsCount-1)*wordIntervalLength, freqMatrix.wordCount, &wg)
	wg.Wait()

	//Combine Tf and Idf to get Tf-Idf
	for i := 0; i < threadsCount-1; i++ {
		startIndex := i * fileIntervalLength
		endIndex := (i + 1) * fileIntervalLength
		wg.Add(1)
		go combineTfAndIidf(freqMatrix, idfs, startIndex, endIndex, &wg)
	}
	wg.Add(1)
	go combineTfAndIidf(freqMatrix, idfs, (threadsCount-1)*fileIntervalLength, len(files), &wg)
	wg.Wait()

	now := time.Now()
	textLine := "Indexing on " + fmt.Sprint(threadsCount) + "threads, take " + fmt.Sprint(now.Sub(start))

	fmt.Println(strings.Repeat("#", len(textLine)))
	fmt.Println(textLine)
	fmt.Println(strings.Repeat("#", len(textLine)))
	fmt.Println()
	return
}
