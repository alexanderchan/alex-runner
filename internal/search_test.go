package runner

import (
	"testing"
	"time"
)

func createTestScoredScripts() []ScoredScript {
	now := time.Now()
	return []ScoredScript{
		{
			Script:        NPMScript{Name: "dev", Command: "next dev"},
			FrecencyScore: 10.0,
			LastUsed:      &now,
			UseCount:      20,
		},
		{
			Script:        NPMScript{Name: "build", Command: "next build"},
			FrecencyScore: 5.0,
			LastUsed:      &now,
			UseCount:      10,
		},
		{
			Script:        NPMScript{Name: "build:cli", Command: "esbuild ./scripts/cli.ts"},
			FrecencyScore: 3.0,
			LastUsed:      &now,
			UseCount:      5,
		},
		{
			Script:        NPMScript{Name: "test", Command: "jest"},
			FrecencyScore: 2.0,
			LastUsed:      &now,
			UseCount:      3,
		},
		{
			Script:        NPMScript{Name: "typecheck", Command: "tsc --noEmit"},
			FrecencyScore: 1.0,
			LastUsed:      &now,
			UseCount:      2,
		},
		{
			Script:        NPMScript{Name: "start-docker:traefik:build", Command: "docker compose down && docker compose -f ./compose.yml build"},
			FrecencyScore: 4.0,
			LastUsed:      &now,
			UseCount:      8,
		},
		{
			Script:        NPMScript{Name: "demo-record", Command: "vhs record > cassette.tape"},
			FrecencyScore: 2.5,
			LastUsed:      &now,
			UseCount:      4,
		},
		{
			Script:        NPMScript{Name: "demo-generate", Command: "vhs cassette.tape"},
			FrecencyScore: 2.0,
			LastUsed:      &now,
			UseCount:      3,
		},
	}
}

func TestSearchScriptsExactMatch(t *testing.T) {
	scripts := createTestScoredScripts()

	results := SearchScripts(scripts, "dev")

	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}

	if results[0].Script.Name != "dev" {
		t.Errorf("expected 'dev' as first result, got '%s'", results[0].Script.Name)
	}
}

func TestSearchScriptsPrefixMatch(t *testing.T) {
	scripts := createTestScoredScripts()

	results := SearchScripts(scripts, "bui")

	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}

	// Both "build" and "build:cli" should match
	foundBuild := false
	foundBuildCli := false

	for _, result := range results {
		if result.Script.Name == "build" {
			foundBuild = true
		}
		if result.Script.Name == "build:cli" {
			foundBuildCli = true
		}
	}

	if !foundBuild {
		t.Error("expected 'build' in results")
	}
	if !foundBuildCli {
		t.Error("expected 'build:cli' in results")
	}

	// "build" should rank higher than "build:cli" because it's a prefix match
	if results[0].Script.Name != "build" {
		t.Errorf("expected 'build' as first result, got '%s'", results[0].Script.Name)
	}
}

func TestSearchScriptsPartialMatch(t *testing.T) {
	scripts := createTestScoredScripts()

	results := SearchScripts(scripts, "type")

	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}

	if results[0].Script.Name != "typecheck" {
		t.Errorf("expected 'typecheck' as first result, got '%s'", results[0].Script.Name)
	}
}

func TestSearchScriptsCommandMatch(t *testing.T) {
	scripts := createTestScoredScripts()

	results := SearchScripts(scripts, "jest")

	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}

	if results[0].Script.Name != "test" {
		t.Errorf("expected 'test' script (contains 'jest' command), got '%s'", results[0].Script.Name)
	}
}

func TestSearchScriptsCaseInsensitive(t *testing.T) {
	scripts := createTestScoredScripts()

	resultsLower := SearchScripts(scripts, "dev")
	resultsUpper := SearchScripts(scripts, "DEV")
	resultsMixed := SearchScripts(scripts, "DeV")

	if len(resultsLower) == 0 || len(resultsUpper) == 0 || len(resultsMixed) == 0 {
		t.Fatal("case insensitive search failed")
	}

	if resultsLower[0].Script.Name != resultsUpper[0].Script.Name ||
		resultsLower[0].Script.Name != resultsMixed[0].Script.Name {
		t.Error("case insensitive search returned different results")
	}
}

func TestSearchScriptsNoMatch(t *testing.T) {
	scripts := createTestScoredScripts()

	results := SearchScripts(scripts, "nonexistent")

	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestSearchScriptsEmptyQuery(t *testing.T) {
	scripts := createTestScoredScripts()

	results := SearchScripts(scripts, "")

	if len(results) != len(scripts) {
		t.Errorf("expected all scripts for empty query, got %d out of %d", len(results), len(scripts))
	}
}

func TestSearchScriptsRankingWithFrecency(t *testing.T) {
	now := time.Now()
	scripts := []ScoredScript{
		{
			Script:       NPMScript{Name: "build", Command: "next build"},
			FrecencyScore: 10.0, // Higher frecency
			LastUsed:     &now,
			UseCount:     20,
		},
		{
			Script:       NPMScript{Name: "build:prod", Command: "NODE_ENV=production next build"},
			FrecencyScore: 2.0, // Lower frecency
			LastUsed:     &now,
			UseCount:     3,
		},
	}

	results := SearchScripts(scripts, "build")

	if len(results) < 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// "build" should rank higher due to exact match
	if results[0].Script.Name != "build" {
		t.Errorf("expected 'build' as first result, got '%s'", results[0].Script.Name)
	}
}

func TestSearchScriptsFuzzyMatch(t *testing.T) {
	scripts := createTestScoredScripts()

	// Fuzzy match should catch typos
	results := SearchScripts(scripts, "buld") // Missing 'i'

	if len(results) == 0 {
		t.Fatal("expected fuzzy match to find 'build'")
	}

	// Should find "build" scripts
	foundBuild := false
	for _, result := range results {
		if result.Script.Name == "build" || result.Script.Name == "build:cli" {
			foundBuild = true
			break
		}
	}

	if !foundBuild {
		t.Error("fuzzy search should find 'build' for query 'buld'")
	}
}

func TestSearchScriptsMultipleMatches(t *testing.T) {
	scripts := createTestScoredScripts()

	results := SearchScripts(scripts, "b")

	// Should match "build" and "build:cli"
	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}

	foundBuild := false
	foundBuildCli := false

	for _, result := range results {
		if result.Script.Name == "build" {
			foundBuild = true
		}
		if result.Script.Name == "build:cli" {
			foundBuildCli = true
		}
	}

	if !foundBuild {
		t.Error("expected 'build' in results")
	}
	if !foundBuildCli {
		t.Error("expected 'build:cli' in results")
	}
}

func TestSearchScriptsWhitespace(t *testing.T) {
	scripts := createTestScoredScripts()

	// Test with leading/trailing whitespace
	results := SearchScripts(scripts, "  dev  ")

	if len(results) == 0 {
		t.Fatal("expected results with whitespace query")
	}

	if results[0].Script.Name != "dev" {
		t.Errorf("expected 'dev' as first result, got '%s'", results[0].Script.Name)
	}
}

func TestSearchScriptsMultiWordQueries(t *testing.T) {
	scripts := createTestScoredScripts()

	tests := []struct {
		name          string
		query         string
		expectMatch   string
		expectResults bool
	}{
		{
			name:          "Multi-word query 'docker build' should match 'start-docker:traefik:build'",
			query:         "docker build",
			expectMatch:   "start-docker:traefik:build",
			expectResults: true,
		},
		{
			name:          "Multi-word query 'build docker' (reversed) should match 'start-docker:traefik:build'",
			query:         "build docker",
			expectMatch:   "start-docker:traefik:build",
			expectResults: true,
		},
		{
			name:          "Multi-word query 'traefik build' should match 'start-docker:traefik:build'",
			query:         "traefik build",
			expectMatch:   "start-docker:traefik:build",
			expectResults: true,
		},
		{
			name:          "Single word 'dmo' should fuzzy match 'demo-record'",
			query:         "dmo",
			expectMatch:   "demo-record",
			expectResults: true,
		},
		{
			name:          "Single word 'demo' should match 'demo-record'",
			query:         "demo",
			expectMatch:   "demo-record",
			expectResults: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := SearchScripts(scripts, tt.query)

			if tt.expectResults {
				if len(results) == 0 {
					t.Errorf("Expected results for query %q but got none", tt.query)
					return
				}

				// Check if the expected match appears in results (preferably first)
				found := false
				position := -1
				for i, result := range results {
					if result.Script.Name == tt.expectMatch {
						found = true
						position = i
						break
					}
				}

				if !found {
					t.Errorf("Expected %q in results for query %q but not found. Got: %v",
						tt.expectMatch, tt.query, getScriptNames(results))
					return
				}

				t.Logf("Query %q matched %q at position %d (total: %d results)",
					tt.query, tt.expectMatch, position, len(results))
			} else {
				if len(results) > 0 {
					t.Errorf("Expected no results for query %q but got %d", tt.query, len(results))
				}
			}
		})
	}
}

// Helper function to extract script names from results for debugging
func getScriptNames(results []ScoredScript) []string {
	names := make([]string, len(results))
	for i, result := range results {
		names[i] = result.Script.Name
	}
	return names
}
