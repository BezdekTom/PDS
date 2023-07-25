package main

import (
	"io/ioutil"
	"log"
	"strings"
)

// Split text on words, ignore given characters
func SplitTextOnWords(text string, textsToIgnore []string) (words []string) {
	text = strings.ToLower(text)
	for _, textToIgnore := range textsToIgnore {
		text = strings.Replace(text, textToIgnore, " ", -1)
	}
	words = strings.Fields(text)
	return
}

// Get names of all files in given folder
func ListDir(dir_name string) (file_names []string) {
	files, err := ioutil.ReadDir(dir_name)
	if err != nil {
		log.Fatal(err)
	}

	file_names = make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			file_names = append(file_names, file.Name())
		}
	}
	return
}
