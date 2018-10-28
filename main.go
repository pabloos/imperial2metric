package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	resultsDir = "results"

	//Regexp
	//terminals
	digit = "[0-9]"

	zero           = "quot"
	quarter        = "#188"
	half           = "#189"
	halfAndQuarter = "#190"

	//non-terminals
	numberRegex = "[1-9]?" + digit + "+"
	mesureRegex = numberRegex + " ?\""

	floatingCommas = "(" + zero + "|" + quarter + "|" + half + "|" + halfAndQuarter + ")"

	quote              = "<ut type=\"unknown\" x=\"" + digit + "+" + "\">&amp;amp;" + floatingCommas + ";</ut>"
	imperialMatchRegex = numberRegex + " ?" + quote + "(" + quote + ")?"
	sourceExp          = "<source.*>.*" + quote + ".*?" + "</source>"

	sourceWithoutQuote = "<source>.*" + numberRegex + ".*</source>"
)

var decimalsMap = map[string]float64{
	zero:           0.0,
	quarter:        0.25,
	half:           0.5,
	halfAndQuarter: 0.75,
}

func getFileContent(path string) (lines string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	return string(content)
}

func imperial2Metric(number float64) float64 {
	if number == 15.5 { //client's exception case
		return 39.5
	}

	return math.Round(number*2.54*10) / 10
}

func getTheNumbers(matches []string, numberRegex string) []float64 {
	rnumber, _ := regexp.Compile(numberRegex)

	numbers := make([]float64, len(matches))

	for i, match := range matches {
		numbers[i], _ = strconv.ParseFloat(rnumber.FindString(match), 64) // https://stackoverflow.com/questions/18951359/how-to-format-floating-point-numbers-into-a-string-using-go?answertab=votes#tab-top
	}

	return getFloatingPoints(matches, numbers)
}

func getFloatingPoints(matches []string, numbers []float64) []float64 {
	var numbersWithDecimals []float64

	decimalRegexp, err := regexp.Compile(floatingCommas)
	if err != nil {
		fmt.Println(err)
	}

	for i, match := range matches {
		decimal := decimalRegexp.FindAllString(match, 1)

		numbersWithDecimals = append(numbersWithDecimals, imperial2Metric(numbers[i]+decimalsMap[decimal[0]]))
	}

	return numbersWithDecimals
}

func getSource(sourceExp string, lines string) string {
	sourceRegexp, err := regexp.Compile(sourceExp)

	if err != nil {
		log.Fatalln("Cannot compile the regular expresion of the source")
	}

	return sourceRegexp.FindAllString(lines, 1)[0] //we are getting the first match from the array
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

/* func convert2Comma(number float64) string {
	//it could may be better to use strconv.FormatFloat instead of fmt.Sprintf
	return strings.Replace(fmt.Sprintf("%.1f", number), ".", ",", 1)
} */

func convert2Comma(number float64) string { // https://godoc.org/golang.org/x/text/message
	return message.NewPrinter(language.Spanish).Sprintf("%.1f", number)
}

func writeSource(firstSource, secondSource, lines string) string {
	regexpSource, _ := regexp.Compile(secondSource)

	return regexpSource.ReplaceAllString(lines, firstSource)
}

func writeOnFile(path, dirname, filename, content string) {
	errMkdir := os.Mkdir(path+"/"+dirname, 0755)
	if errMkdir != nil {
		fmt.Println(errMkdir)
	}

	err := ioutil.WriteFile(path+"/"+dirname+"/"+filename, []byte(content), 0666)
	if err != nil {
		fmt.Println(err)
	}
}

func transformFile(file string) {

	/*
		we have a function that can write in source and target and another for write only in source
		so in order to write in targe we need to:
		1) get the first source
		2) write in source and target
		3) replace the new second source with the first one
	*/

	lines := getFileContent(file)

	firstSource := getSource(sourceExp, lines)

	matches := getImperialMatches(imperialMatchRegex, lines)

	numbers := getTheNumbers(matches, numberRegex)

	lines = writeSourceAndTarget(matches, numbers, lines)

	secondSource := getSource(sourceWithoutQuote, lines)

	lines = writeSource(firstSource, secondSource, lines)

	writeOnFile(filepath.Dir(file), resultsDir, filepath.Base(file), lines)
}

func main() {

	files := os.Args[1:]

	if len(files) < 1 {
		fmt.Print("\nYou have not introduced any file/s\n\n")
		fmt.Print("Usage: \n --> $ imperial2metric files...\n\n")

		os.Exit(1)
	}

	for _, file := range files {
		transformFile(file)
	}
}
