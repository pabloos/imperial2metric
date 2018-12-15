package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const resultsDir = "results"

func getSource(sourceExp string, lines string) (string, error) {
	sourceRegexp, err := regexp.Compile(sourceExp)

	if err != nil {
		return "", err
	}

	if len(sourceRegexp.FindAllString(lines, 1)) < 1 {
		return "", errors.New("seems to have not a <source> tag. Skipping his conversion")
	}

	return sourceRegexp.FindAllString(lines, 1)[0], nil //we are getting the first match from the array
}

func getImperialMatches(imperialMatchRegex string, lines string) []string {
	r, _ := regexp.Compile(imperialMatchRegex)

	return r.FindAllString(lines, 3)
}

func writeSourceAndTarget(matches []string, numbers []float64, lines string) string {
	for i, match := range matches {
		r, _ := regexp.Compile(match)

		// lines = r.ReplaceAllString(lines, fmt.Sprintf("%.1f", numbers[i])) //cambia source y target

		lines = r.ReplaceAllString(lines, convert2Comma(numbers[i])) //cambia source y target
	}

	return lines
}

func getTheNumbers(matches []string, numberRegex string) []float64 {
	rnumber, _ := regexp.Compile(numberRegex)

	numbers := make([]float64, len(matches))

	for i, match := range matches {
		numbers[i], _ = strconv.ParseFloat(rnumber.FindString(match), 64) // https://stackoverflow.com/questions/18951359/how-to-format-floating-point-numbers-into-a-string-using-go?answertab=votes#tab-top
	}

	return getFloatingPoints(matches, numbers)
}

func writeSource(firstSource, secondSource, lines string) string {
	regexpSource, _ := regexp.Compile(secondSource)

	return regexpSource.ReplaceAllString(lines, firstSource)
}

func WriteOnFile(path, dirname, filename string, bytes []byte) {
	err := ioutil.WriteFile(path+"/"+dirname+"/"+filename, bytes, 0666)
	if err != nil {
		fmt.Println(err)
	}
}

func TransformFile(input io.Reader) (io.Reader, error) {

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

	lines := string(buff.Bytes()) // [:] ?

	firstSource, err := getSource(sourceExp, lines)
	if err != nil {
		return strings.NewReader(lines), nil
	}

	matches := getImperialMatches(imperialMatchRegex, lines)

	numbers := getTheNumbers(matches, numberRegex)

	lines = writeSourceAndTarget(matches, numbers, lines)

	secondSource, err := getSource(sourceWithoutQuote, lines)
	if err != nil {
		return input, err
	}

	lines = writeSource(firstSource, secondSource, lines)

	return strings.NewReader(lines), nil
}
