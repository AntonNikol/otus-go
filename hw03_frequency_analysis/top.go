package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	Quantity int
	Text     string
}

// Top10 Возвращает слайс с 10-ю наиболее часто встречаемыми в тексте словами.
func Top10(str string) []string {
	if len(str) == 0 {
		return []string{}
	}

	words := countWords(str)
	sortedWords := sortWords(words)

	if len(sortedWords) > 10 {
		sortedWords = sortedWords[:10]
	}
	return extractText(sortedWords)
}

func sortWords(words map[string]int) []Word {
	s := make([]Word, 0, len(words))
	for k, v := range words {
		s = append(s, Word{Text: k, Quantity: v})
	}

	sort.Slice(s, func(i, j int) bool {
		if s[i].Quantity == s[j].Quantity {
			return s[i].Text < s[j].Text
		}
		return s[i].Quantity > s[j].Quantity
	})

	return s
}

func countWords(str string) map[string]int {
	m := map[string]int{}
	for _, v := range strings.Fields(str) {
		if v != "" {
			m[v]++
		}
	}
	return m
}

func extractText(words []Word) []string {
	result := make([]string, 0, len(words))
	for _, v := range words {
		result = append(result, v.Text)
	}
	return result
}
