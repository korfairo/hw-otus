package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func Top10(s string) []string {
	words := strings.Fields(s)

	wordsToCountMap := countFrequency(words)

	uniqueWords := mapKeysToSlice(wordsToCountMap)

	sort.SliceStable(uniqueWords, func(i, j int) bool {
		if wordsToCountMap[uniqueWords[i]] != wordsToCountMap[uniqueWords[j]] {
			return wordsToCountMap[uniqueWords[i]] > wordsToCountMap[uniqueWords[j]]
		}
		return uniqueWords[i] < uniqueWords[j]
	})

	if len(uniqueWords) > 10 {
		return uniqueWords[:10]
	}
	return uniqueWords
}

// countFrequency returns map with words as keys and their counts as values.
func countFrequency(words []string) map[string]int {
	wordsToCount := make(map[string]int, len(words))
	for _, word := range words {
		word = trimEdgePunctuationMarks(word)
		if word == "" {
			continue
		}
		word = strings.ToLower(word)
		wordsToCount[word]++
	}
	return wordsToCount
}

// trimEdgePunctuationMarks trims punctuation marks in the first and the last position in the string.
func trimEdgePunctuationMarks(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	lastRuneIdx := len(runes) - 1

	if unicode.IsPunct(runes[0]) {
		if lastRuneIdx == 0 {
			return ""
		}
		runes = runes[1:]
		lastRuneIdx--
	}

	if unicode.IsPunct(runes[lastRuneIdx]) {
		runes = runes[:lastRuneIdx]
	}

	return string(runes)
}

func mapKeysToSlice[T interface{}](m map[string]T) []string {
	slice := make([]string, 0, len(m))
	for key := range m {
		slice = append(slice, key)
	}
	return slice
}
