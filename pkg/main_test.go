package pkg

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestImperial2Metric(t *testing.T) {
	testTable := []struct {
		value    float64
		expected float64
	}{
		{6, 15.2},
		{24, 61},
		{9, 22.9},
		{9.25, 23.5},
	}

	for _, test := range testTable {
		t.Run(strconv.FormatFloat(test.value, 2, 'f', 64), func(t *testing.T) {
			if imperial2Metric(test.value) != test.expected {
				t.Errorf("%f should be %f", test.value, test.expected)
			}
		})
	}
}

func TestConvert2Comma(t *testing.T) {
	testTable := []struct {
		source   float64
		expected string
	}{
		{
			1.0,
			"1,0",
		},
		{
			10.3,
			"10,3",
		},
	}

	for _, test := range testTable {
		t.Run(fmt.Sprintf("Test case %f", test.source), func(t *testing.T) {
			result := convert2Comma(test.source)
			if result != test.expected {
				t.Errorf(fmt.Sprintf("Test case %f failed. Got %s and expected %s", test.source, result, test.expected))
			}
		})
	}
}

func TestGetTheNumbers(t *testing.T) {
	testTable := []struct {
		matches  []string
		expected []float64
	}{
		{
			matches: []string{
				"6<ut type=\"unknown\" x=\"1\">&amp;amp;quot;</ut>",
				"24<ut type=\"unknown\" x=\"2\">&amp;amp;quot;</ut>",
				"9<ut type=\"unknown\" x=\"3\">&amp;amp;#188;</ut><ut type=\"unknown\" x=\"4\">&amp;amp;quot;</ut>",
			},
			expected: []float64{
				15.2,
				61.0,
				23.5,
			},
		},
		{
			matches:  []string{},
			expected: []float64{},
		},
	}

	for i, test := range testTable {
		t.Run(fmt.Sprintf("Test case: %d", i), func(t *testing.T) {
			for j, number := range getTheNumbers(test.matches, numberRegex) {
				if number != test.expected[j] {
					t.Errorf("Test %d failed", i)
				}
			}
		})
	}
}

func TestGetSource(t *testing.T) {
	testTable := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "",
			source:   "<skeleton>&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;</skeleton><translatable blockId=\"5\" datatype=\"htmls\"><segment segmentId=\"1\" modified=\"true\"><source>Size 10 has a 6<ut type=\"unknown\" x=\"1\">&amp;amp;quot;</ut> rise, a 24<ut type=\"unknown\" x=\"2\">&amp;amp;quot;</ut> inseam, and a 9<ut type=\"unknown\" x=\"3\">&amp;amp;#188;</ut><ut type=\"unknown\" x=\"4\">&amp;amp;quot;</ut> leg opening.</source><target creationid=\"jsalmeron\" creationdate=\"20181016T082326Z\" score=\"100\" seginfo=\"&lt;root username=&quot;jsalmeron&quot; timestamp=&quot;20181016T082326Z&quot; tool-name=&quot;TSO&quot; tool-version=&quot;5.5.0&quot; set=&quot;3B1&quot; phase-name=&quot;Translation-1&quot;/&gt;\">La talla M tiene 15,2 de tiro, 61,0 de entrepierna y 23,5 de abertura de pernera.</target><comments><comment creationid=\"alejandroruiz\" creationdate=\"20181016T125347Z\" type=\"text\">Size \"10\" in source but \"M\" in target, as handled by brand captain previously. Please check if ok.</comment></comments><revisions><revision creationid=\"ProjectDirector\" creationdate=\"20180921T195408Z\" type=\"target\" seginfo=\"&lt;root phase-name=&quot;pd-1&quot;/&gt;\" id=\"1\"><target score=\"-1\">La talla M tiene 15,2 cm de tiro, 60,9 cm de entrepierna y 22,8 cm de apertura de piernas.</target></revision></revisions></segment></translatable>",
			expected: "<source>Size 10 has a 6<ut type=\"unknown\" x=\"1\">&amp;amp;quot;</ut> rise, a 24<ut type=\"unknown\" x=\"2\">&amp;amp;quot;</ut> inseam, and a 9<ut type=\"unknown\" x=\"3\">&amp;amp;#188;</ut><ut type=\"unknown\" x=\"4\">&amp;amp;quot;</ut> leg opening.</source>",
		},
	}

	for _, test := range testTable {
		t.Run(fmt.Sprintf("Test case %s", test.name), func(t *testing.T) {
			result, err := getSource(sourceExp, test.source)
			if err != nil {
				t.Errorf("Error on test %s: %v", test.name, err)
			}

			if result != test.expected {
				t.Errorf("Test case: %s failed: expected %s, got %s", test.name, test.expected, result)
			}
		})
	}
}

func TestGetImperialMatches(t *testing.T) {
	lines, _ := ioutil.ReadFile("CAT_462457_masterCatalog_RalphLauren_EU-product_es--es-ES-1-XML_DMW-es-ES#CNUUK#.xml-TXML-es-ES#PQ_HNJNR#.txml")

	testTable := []struct {
		name     string
		source   string
		expected []string
	}{
		{
			name:   "",
			source: string(lines),
			expected: []string{
				"6<ut type=\"unknown\" x=\"1\">&amp;amp;quot;</ut>",
				"24<ut type=\"unknown\" x=\"2\">&amp;amp;quot;</ut>",
				"9<ut type=\"unknown\" x=\"3\">&amp;amp;#188;</ut><ut type=\"unknown\" x=\"4\">&amp;amp;quot;</ut>",
			},
		},
	}

	for _, test := range testTable {
		t.Run("", func(t *testing.T) {
			result := getImperialMatches(imperialMatchRegex, test.source)
			for i, matchResult := range result {
				if matchResult != test.expected[i] {
					t.Errorf(fmt.Sprintf("Test case %s failed: expected %s and got %s", test.name, test.expected[i], result))
				}
			}
		})
	}
}

func TestGetFloatingPoints(t *testing.T) {
	testTable := []struct {
		name     string
		matches  []string
		numbers  []float64
		expected []float64
	}{
		{
			name: "",
			matches: []string{
				"6<ut type=\"unknown\" x=\"1\">&amp;amp;quot;</ut>",
				"24<ut type=\"unknown\" x=\"2\">&amp;amp;quot;</ut>",
				"9<ut type=\"unknown\" x=\"3\">&amp;amp;#188;</ut><ut type=\"unknown\" x=\"4\">&amp;amp;quot;</ut>",
			},
			numbers: []float64{
				6.0, 24.0, 9.0,
			},
			expected: []float64{
				15.2, 61.0, 23.5,
			},
		},
	}

	for _, test := range testTable {
		t.Run("", func(t *testing.T) {
			result := getFloatingPoints(test.matches, test.numbers)

			for i, number := range result {
				if number != test.expected[i] {
					t.Errorf("Test %s failed. Expected %f and got %f", test.name, test.expected[i], number)
				}
			}
		})
	}
}
