package search

import (
	"sort"
	"strings"

	"github.com/gjtiquia/ifg/internal/config"
)

type scoredEntry struct {
	entry config.Entry
	score int
}

func Match(entries []config.Entry, query string) []config.Entry {
	if query == "" {
		return entries
	}

	tokens := tokenize(query)
	var scored []scoredEntry

	for _, entry := range entries {
		score := matchEntry(entry, tokens)
		if score > 0 {
			scored = append(scored, scoredEntry{entry: entry, score: score})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	result := make([]config.Entry, len(scored))
	for i, se := range scored {
		result[i] = se.entry
	}

	return result
}

func tokenize(query string) []string {
	var tokens []string
	for _, token := range strings.Fields(query) {
		if token != "" {
			tokens = append(tokens, strings.ToLower(token))
		}
	}
	return tokens
}

func matchEntry(entry config.Entry, tokens []string) int {
	allMatch := true
	minScore := 0

	for _, token := range tokens {
		score := matchToken(entry, token)
		if score == 0 {
			allMatch = false
			break
		}
		if minScore == 0 || score < minScore {
			minScore = score
		}
	}

	if !allMatch {
		return 0
	}

	return minScore
}

func matchToken(entry config.Entry, token string) int {
	lowerCommand := strings.ToLower(entry.Command)
	lowerTitle := strings.ToLower(entry.Title)

	if strings.Contains(lowerCommand, token) {
		return 100
	}

	if strings.Contains(lowerTitle, token) {
		return 50
	}

	for _, desc := range entry.Description {
		if strings.Contains(strings.ToLower(desc), token) {
			return 25
		}
	}

	return 0
}
