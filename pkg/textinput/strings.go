package textinput

import (
	"strings"
)

// Split slices s into all substrings separated by any of given seps and returns a slice of
// the substrings between those separators with quote awareness.
// Examples:
//
//   - create_question 1 "How many letters are there in the English alphabet?” 26
//     > []string{"create_question", "1", "\"How many letters are there in the English alphabet?\"", "26"}
func Split(s string, seps ...rune) []string {
	return SplitWithOptions(s, seps, true, false)
}

// SplitWithOptions slices s into all substring separated by any if given seps and return a slice of
// the substrings between those separators with options of quote awareness and/or to include seps in the results.
// Examples:
//
//  1. Without Quote Awareness (seps: []rune{" "})
//     - create_question 1 "Who am I?" me
//     > []string{"create_question", "1", "\"Who", "am", "I?\"", "me"}
//  2. With Quote Awareness (seps: []rune{" "})
//     - create_question 1 "Who am I?” me
//     > []string{"create_question", "1", "\"Who am I?\”", "me"}
//  3. With Quote Awareness incude Separators (seps: []rune{" "})
//     - create_question 1 "Who am I?” me
//     > []string{"create_question", " ", "1", " ", "\"Who am I?\”", " ", "me"}
//  4. Without Quote Awareness incude Separators (multiple) (seps: []rune{' ', ',', '?', '(', ')'})
//     - create_question 1 "1, 2 or (3)?" 3
//     > []string{"create_question", " ", "1", " ", "\"1", ",", " ", "2", " ", "or",
//     " ", "(", "3", ")", "?", "\"", " ", "3"}
func SplitWithOptions(s string, seps []rune, quoteAware, includeSep bool) []string {
	var (
		quoteFlag bool
		stream    []rune
		subseps   []string
		msep      = make(map[rune]struct{})
	)

	for _, sep := range seps {
		msep[sep] = struct{}{}
	}

	for _, c := range s {
		if quoteAware && c == '"' {
			quoteFlag = !quoteFlag
		}

		if _, ok := msep[c]; ok && !quoteFlag {
			if len(stream) != 0 {
				subseps = append(subseps, string(stream))
			}

			if includeSep {
				subseps = append(subseps, string(c))
			}

			stream = []rune{}
			continue
		}

		stream = append(stream, c)
	}

	subseps = append(subseps, string(stream))

	return subseps
}

var stringsNumber = map[string]string{
	"zero":  "0",
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
	"ten":   "10",
}

// RecognizedAsNumber converts all substrings in the given s that represent numbers.
// Limitation only handle 0 to 10.
// Examples:
//
//   - one, two, three or four?! -> 1, 2, 3 or 4?!
//   - loss (one) usd -> loss (1) usd
func RecognizedAsNumber(s string) string {
	var (
		sep                      rune = ' '
		commonSymbolsCoverNumber      = []rune{',', '.', '?', '!', '(', ')'}
	)

	parts := SplitWithOptions(s, append(commonSymbolsCoverNumber, sep), false, true)
	for i, part := range parts {
		part := strings.ToLower(part)
		v, ok := stringsNumber[part]
		if !ok {
			continue
		}
		parts[i] = v
	}

	return strings.Join(parts, "")
}
