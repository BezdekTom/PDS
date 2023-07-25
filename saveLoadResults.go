package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

/*
##########################################
Save tf-idf results to given folder
##########################################
*/

//Safe results of tf-idf to file
func SaveResults(indexMap *SafeMap, frequencyMatrix *FrequencyMatrix, textFiles []string, folderPath string) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	mapPath := path.Join(folderPath, "indexMap.txt")
	indexMap.saveToFile(mapPath)
	matrixPath := path.Join(folderPath, "frequencyMatrix.txt")
	frequencyMatrix.saveToFile(matrixPath)
	textFilesPath := path.Join(folderPath, "textFiles.txt")
	saveFileIndexes(textFiles, textFilesPath)
}

//save frequation matrix to given text file
func (safeMap *SafeMap) saveToFile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for w, i := range safeMap.items {
		i_string := strconv.Itoa(i)
		_, err := f.WriteString(w + "\t" + i_string + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

//Save frequention matrix to given text file
func (frequencyMatrix *FrequencyMatrix) saveToFile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	numberOfFiles := strconv.Itoa(len(frequencyMatrix.matrix))
	numberOfWords := strconv.Itoa(frequencyMatrix.wordCount)
	_, err = f.WriteString(numberOfFiles + "\t" + numberOfWords + "\n")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(frequencyMatrix.matrix); i++ {
		for j := 0; j < frequencyMatrix.wordCount; j++ {
			number := fmt.Sprint(frequencyMatrix.matrix[i][j])
			_, err := f.WriteString(number + "\t")
			if err != nil {
				log.Fatal(err)
			}
		}
		_, err := f.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

//Save file names and coresponding indexes in frequention matrix to given text file
func saveFileIndexes(textFiles []string, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lenString := strconv.Itoa(len(textFiles))
	_, err = f.WriteString(lenString + "\n")
	if err != nil {
		log.Fatal(err)
	}

	for idx, fileName := range textFiles {
		idxString := strconv.Itoa(idx)
		_, err := f.WriteString(idxString + "\t" + fileName + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

/*
##########################################
Load tf-idf results from given folder
##########################################
*/

//Load tf-idf resluts from text files in given folder
func LoadTextIndex(folderPath string) (indexMap *SafeMap, frequencyMatrix *FrequencyMatrix, fileIndexes []string) {
	mapPath := path.Join(folderPath, "indexMap.txt")
	indexMap = loadIndexMap(mapPath)
	matrixPath := path.Join(folderPath, "frequencyMatrix.txt")
	frequencyMatrix = loadFreqMatrix(matrixPath)
	textFilesPath := path.Join(folderPath, "textFiles.txt")
	fileIndexes = loadFileIndexes(textFilesPath)
	return
}

//Load index map from given text file
func loadIndexMap(filePath string) (indexMap *SafeMap) {
	indexMap = &SafeMap{
		indexCount: 0,
		items:      map[string]int{},
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		idx, err := strconv.Atoi(words[1])
		if err != nil {
			panic(err)
		}
		indexMap.items[words[0]] = idx
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

//Load frequnetion matrix from given text file
func loadFreqMatrix(filePath string) (freqMatrix *FrequencyMatrix) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	fileCountS := scanner.Text()
	scanner.Scan()
	wordCountS := scanner.Text()
	filesCount, err := strconv.Atoi(fileCountS)
	if err != nil {
		panic(err)
	}

	wordsCount, err := strconv.Atoi(wordCountS)
	if err != nil {
		panic(err)
	}

	freqMatrix = &FrequencyMatrix{
		wordCount: wordsCount,
		capacity:  wordsCount,
		matrix:    make([][]float64, filesCount),
	}
	freqMatrix.Init()

	for fileIdx := 0; fileIdx < filesCount; fileIdx++ {
		for wordIdx := 0; wordIdx < freqMatrix.wordCount; wordIdx++ {
			scanner.Scan()
			countS := scanner.Text()
			countF, err := strconv.ParseFloat(strings.TrimSpace(countS), 64)
			if err != nil {
				panic(err)
			}
			freqMatrix.matrix[fileIdx][wordIdx] = countF
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

//Load file names and coresponding indexes in frequention matrix from given text file
func loadFileIndexes(filePath string) (textFiles []string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()
	words := strings.Fields(firstLine)
	filesCount, err := strconv.Atoi(words[0])
	if err != nil {
		panic(err)
	}
	textFiles = make([]string, filesCount)

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		index, err := strconv.Atoi(words[0])
		if err != nil {
			panic(err)
		}
		textFiles[index] = words[1]

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
