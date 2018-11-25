package pkg

import (
	"log"
	"regexp"
	"strconv"
)

func getSource(sourceExp string, lines string) string {
	sourceRegexp, err := regexp.Compile(sourceExp)

	if err != nil {
		log.Fatalln("Cannot compile the regular expression of the source")
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
