package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

const resultsDir = "results"

func WriteOnFile(path, dirname, filename string, bytes []byte) {
	err := ioutil.WriteFile(path+"/"+dirname+"/"+filename, bytes, 0666)
	if err != nil {
		fmt.Println(err)
	}
}

func TransformFile(input io.Reader) io.Reader {

	/*
		we have a function that can write in source and target and another for write only in source
		so in order to write in targe we need to:
		1) get the first source
		2) write in source and target
		3) replace the new second source with the first one
	*/

	buff := new(bytes.Buffer)

	_, err := buff.ReadFrom(input)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	lines := string(buff.Bytes())

	firstSource := getSource(sourceExp, lines)

	matches := getImperialMatches(imperialMatchRegex, lines)

	numbers := getTheNumbers(matches, numberRegex)

	lines = writeSourceAndTarget(matches, numbers, lines)

	secondSource := getSource(sourceWithoutQuote, lines)

	lines = writeSource(firstSource, secondSource, lines)

	return strings.NewReader(lines)
}
