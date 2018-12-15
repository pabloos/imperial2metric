package pkg

import (
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
