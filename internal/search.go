package runner

import (
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/schollz/closestmatch"
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

	// Split query into words to detect multi-word queries
	queryWords := strings.Fields(query)
	isMultiWord := len(queryWords) > 1

	// For multi-word queries, use closestmatch for better results
	if isMultiWord {
		return searchWithClosestMatch(scoredScripts, query)
	}

	// Single word query - use existing fuzzy logic
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

// searchWithClosestMatch uses bag-of-words matching for multi-word queries
// This handles queries like "build docker" or "docker build" effectively
func searchWithClosestMatch(scoredScripts []ScoredScript, query string) []ScoredScript {
	if len(scoredScripts) == 0 {
		return scoredScripts
	}

	// Build lists of script names and commands for matching
	scriptNames := make([]string, len(scoredScripts))
	scriptCommands := make([]string, len(scoredScripts))
	scriptCombined := make([]string, len(scoredScripts))

	for i, scored := range scoredScripts {
		scriptNames[i] = strings.ToLower(scored.Script.Name)
		scriptCommands[i] = strings.ToLower(scored.Script.Command)
		scriptCombined[i] = scriptNames[i] + " " + scriptCommands[i]
	}

	// Create closestmatch instances with n-gram sizes optimized for script names
	cmNames := closestmatch.New(scriptNames, []int{2, 3, 4})
	cmCommands := closestmatch.New(scriptCommands, []int{2, 3, 4})
	cmCombined := closestmatch.New(scriptCombined, []int{2, 3, 4})

	// Get matches from each matcher
	nameMatches := cmNames.ClosestN(query, len(scoredScripts))
	commandMatches := cmCommands.ClosestN(query, len(scoredScripts))
	combinedMatches := cmCombined.ClosestN(query, len(scoredScripts))

	// Build a rank map: script -> rank
	rankMap := make(map[string]int)

	// Rank based on position in match lists (earlier = better)
	// Combined matches get highest weight, then name, then command
	for i, match := range combinedMatches {
		if match != "" {
			// Extract name from combined (before the space)
			name := strings.SplitN(match, " ", 2)[0]
			if _, exists := rankMap[name]; !exists {
				rankMap[name] = 1000 - (i * 10) // Top match gets 1000, decreasing by 10
			}
		}
	}

	for i, match := range nameMatches {
		if match != "" {
			if _, exists := rankMap[match]; !exists {
				rankMap[match] = 800 - (i * 10)
			}
		}
	}

	for i, cmdMatch := range commandMatches {
		if cmdMatch != "" {
			// Find script by command
			for j, cmd := range scriptCommands {
				if cmd == cmdMatch {
					name := scriptNames[j]
					if _, exists := rankMap[name]; !exists {
						rankMap[name] = 600 - (i * 10)
					}
					break
				}
			}
		}
	}

	// Build results with ranks
	var results []searchResult
	for _, scored := range scoredScripts {
		name := strings.ToLower(scored.Script.Name)
		if rank, hasRank := rankMap[name]; hasRank && rank > 0 {
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
