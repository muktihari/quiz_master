package textinput_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/muktihari/quiz_master/pkg/textinput"
)

func TestSplit(t *testing.T) {
	tt := []struct {
		Name     string
		Input    string
		Expected []string
	}{
		{
			Name:     "valid string without quote",
			Input:    "delete_question 1",
			Expected: []string{"delete_question", "1"},
		},
		{
			Name:  "valid string with quote",
			Input: "create_question 1 \"How many letters are there in the English alphabet?\"",
			Expected: []string{"create_question", "1",
				"\"How many letters are there in the English alphabet?\""},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			res := textinput.Split(tc.Input, ' ')
			if diff := cmp.Diff(res, tc.Expected); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSplitWithOptions(t *testing.T) {
	tt := []struct {
		Name       string
		Input      string
		Seps       []rune
		QuoteAware bool
		IncludeSep bool
		Expected   []string
	}{
		{
			Name:       "valid string with multiple separators and include separators, quote aware is false",
			Input:      "create_question 1 \"1, 2 or (3)?\" 3",
			Seps:       []rune{' ', ',', '?', '(', ')'},
			QuoteAware: false,
			IncludeSep: true,
			Expected: []string{"create_question", " ", "1", " ", "\"1", ",", " ", "2", " ",
				"or", " ", "(", "3", ")", "?", "\"", " ", "3"},
		},
		{
			Name:       "valid string with multiple separators and include separators, quote aware is true",
			Input:      "create_question 1 \"1, 2 or (3)?\" 3",
			Seps:       []rune{' ', ',', '?', '(', ')'},
			QuoteAware: true,
			IncludeSep: true,
			Expected:   []string{"create_question", " ", "1", " ", "\"1, 2 or (3)?\"", " ", "3"},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			parts := textinput.SplitWithOptions(tc.Input, tc.Seps, tc.QuoteAware, tc.IncludeSep)
			if diff := cmp.Diff(tc.Expected, parts); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestRecognizedAsNumber(t *testing.T) {
	tt := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "string without number",
			Input:    "cat",
			Expected: "cat",
		},
		{
			Name:     "string with valid number",
			Input:    "zero, one!, two, three, four, five, six, seven, eight or (nine)",
			Expected: "0, 1!, 2, 3, 4, 5, 6, 7, 8 or (9)",
		},
		{
			Name:     "string contains number but not counted as number representative",
			Input:    "walking2",
			Expected: "walking2",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			s := textinput.RecognizedAsNumber(tc.Input)
			if s != tc.Expected {
				t.Fatalf("expected: %s, got: %s", tc.Expected, s)
			}
		})
	}
}
