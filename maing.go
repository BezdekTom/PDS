package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
)

var threadProfile = pprof.Lookup("threadcreate")
var textFolder = "/home/tom/go/src/PDS/PIP0015_projekt_PDS/texts"
var resultFolder = "/home/tom/go/src/PDS/PIP0015_projekt_PDS/results"

// Return text files ordered by relevanced based on given text
func searchInTextFiles(textToFind string, indexFolder string) (filesNames []string) {
	idxMap, freqMatrix, files := LoadTextIndex(indexFolder)
	textsToIgnore := []string{",", ".", ";", ")", "(", "]", "[", "\"", "!", "?", "''"}
	words := SplitTextOnWords(textToFind, textsToIgnore)
	filesNames = SearchText(words, freqMatrix, idxMap, files)
	return
}

// Print result of search
func printSearchResults(filesByRelevance []string, numberOfResults int) {
	if numberOfResults > len(filesByRelevance) {
		numberOfResults = len(filesByRelevance)
	}
	for i := 0; i < numberOfResults; i++ {
		fmt.Println(i+1, ") ", filesByRelevance[i])
	}
}

// Order text files by relevance and print them
func SearchTextFile(textToSearch string, resultFolder string, numberOfResults int) {
	resultFiles := searchInTextFiles(textToSearch, resultFolder)
	printSearchResults(resultFiles, numberOfResults)
}

// Run tf-idf on all text files in given folder and save the results
func RunTfIdf(textFolder string, resultFolder string, thredCount int) {
	idxMap, freqMatrix, textFiles := ComputeTfIdf(textFolder, thredCount)
	SaveResults(idxMap, freqMatrix, textFiles, resultFolder)
}

//Interactive search console app
func BookSearch(resultFolder string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println()
		fmt.Print("What do you want to search (!q for quit): ")
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		input := scanner.Text()
		if input == "!q" {
			return
		}
		SearchTextFile(input, resultFolder, 5)
	}
}

func main() {
	var numOfThreads int = threadProfile.Count()
	var err error
	args := os.Args
	if len(args) >= 2 {
		numOfThreads, err = strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
	}

	RunTfIdf(textFolder, resultFolder, numOfThreads)
	BookSearch(resultFolder)
}
