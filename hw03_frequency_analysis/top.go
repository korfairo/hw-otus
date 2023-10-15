package hw03frequencyanalysis

import (
	"errors"
	"sort"
	"strings"
)

var (
	ErrEmptyString    = errors.New("the string is empty")
	ErrTooFewWordsErr = errors.New("the string contains less than 10 words separated by space")
)

var punctuationMarks = []rune{',', '.', '!', '-'}

func Top10(fuzzySearch bool, s string) ([]string, error) {
	if s == "" {
		return []string{}, ErrEmptyString
	}

	words := strings.Fields(s)

	if len(words) < 10 {
		return []string{}, ErrTooFewWordsErr
	}

	wordsToCountMap := countFrequency(fuzzySearch, words)

	uniqueWords := mapKeysToSlice(wordsToCountMap)
	sort.SliceStable(uniqueWords, func(i, j int) bool {
		if wordsToCountMap[uniqueWords[i]] != wordsToCountMap[uniqueWords[j]] {
			return wordsToCountMap[uniqueWords[i]] > wordsToCountMap[uniqueWords[j]]
		}
		return uniqueWords[i] < uniqueWords[j]
	})

	return uniqueWords[:10], nil
}

// countFrequency returns map with words as keys and their counts as values.
func countFrequency(fuzzySearch bool, words []string) map[string]int {
	wordsToCount := make(map[string]int, len(words))
	for _, word := range words {
		if fuzzySearch {
			word = trimEdgePunctuationMarks(word)

			if word == "" {
				continue
			}

			word = strings.ToLower(word)
		}

		wordsToCount[word]++
	}
	return wordsToCount
}

// isPunctuationMark checks if the rune is a punctuation mark, declared in punctuationMarks global variable.
func isPunctuationMark(r rune) bool {
	for _, mark := range punctuationMarks {
		if mark == r {
			return true
		}
	}
	return false
}

// trimEdgePunctuationMarks trims punctuation marks in the first and the last position in the string.
func trimEdgePunctuationMarks(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	lastRuneIdx := len(runes) - 1

	if isPunctuationMark(runes[0]) {
		if lastRuneIdx == 0 {
			return ""
		}

		runes = runes[1:]
		lastRuneIdx--
	}

	if isPunctuationMark(runes[lastRuneIdx]) {
		runes = runes[:lastRuneIdx]
	}

	return string(runes)
}

func mapKeysToSlice[T interface{}](m map[string]T) []string {
	slice := make([]string, 0, len(m))
	for key, _ := range m {
		slice = append(slice, key)
	}
	return slice
}
