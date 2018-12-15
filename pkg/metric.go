package pkg

import (
	"fmt"
	"math"
	"regexp"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var decimalsMap = map[string]float64{
	zero:           0.0,
	quarter:        0.25,
	half:           0.5,
	halfAndQuarter: 0.75,
}

func imperial2Metric(number float64) float64 {
	if number == 15.5 { //client's exception case
		return 39.5
	}
	return math.Round(number*2.54*10) / 10
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

func convert2Comma(number float64) string { // https://godoc.org/golang.org/x/text/message
	return message.NewPrinter(language.Spanish).Sprintf("%.1f", number)
}
