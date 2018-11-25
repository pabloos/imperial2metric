package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

const resultsDir = "results"

func getFileContent(content []byte) (lines string) {
	return string(content)
}

func writeOnFile(path, dirname, filename, content string) {
	err := ioutil.WriteFile(path+"/"+dirname+"/"+filename, []byte(content), 0666)
	if err != nil {
		fmt.Println(err)
	}
}

func TransformFile(input io.Reader, file string) string {

	/*
		we have a function that can write in source and target and another for write only in source
		so in order to write in targe we need to:
		1) get the first source
		2) write in source and target
		3) replace the new second source with the first one
	*/

	/* 	content, err := ioutil.ReadFile(file)
	   	if err != nil {
	   		log.Fatalf("readLines: %s", err)
	   	} */

	buff := new(bytes.Buffer)

	_, err := buff.ReadFrom(input)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	lines := getFileContent(buff.Bytes())

	firstSource := getSource(sourceExp, lines)

	matches := getImperialMatches(imperialMatchRegex, lines)

	numbers := getTheNumbers(matches, numberRegex)

	lines = writeSourceAndTarget(matches, numbers, lines)

	secondSource := getSource(sourceWithoutQuote, lines)

	lines = writeSource(firstSource, secondSource, lines)

	return lines
}
