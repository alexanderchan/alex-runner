package runner

import (
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type searchResult struct {
	scored ScoredScript
	rank   int
}

func SearchScripts(scoredScripts []ScoredScript, query string) []ScoredScript {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return scoredScripts
	}

	var results []searchResult

	for _, scored := range scoredScripts {
		scriptName := strings.ToLower(scored.Script.Name)
		scriptCommand := strings.ToLower(scored.Script.Command)

		rank := 0

		// Exact match on name gets highest priority
		if scriptName == query {
			rank = 1000
		} else if strings.HasPrefix(scriptName, query) {
			rank = 500
		} else if strings.Contains(scriptName, query) {
			rank = 300
		} else if fuzzy.Match(query, scriptName) {
			// Fuzzy match on name
			rank = 200
		} else if strings.Contains(scriptCommand, query) {
			// Match in command gets lower priority
			rank = 100
		} else if fuzzy.Match(query, scriptCommand) {
			// Fuzzy match in command
			rank = 50
		}

		if rank > 0 {
			results = append(results, searchResult{
				scored: scored,
				rank:   rank,
			})
		}
	}

	// Sort by rank descending, then by frecency score descending
	sort.Slice(results, func(i, j int) bool {
		if results[i].rank != results[j].rank {
			return results[i].rank > results[j].rank
		}
		return results[i].scored.FrecencyScore > results[j].scored.FrecencyScore
	})

	// Extract sorted scripts
	var searchedScripts []ScoredScript
	for _, result := range results {
		searchedScripts = append(searchedScripts, result.scored)
	}

	return searchedScripts
}
