package runner

import (
	"testing"
	"time"
)

func TestCalculateTimeScore(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected float64
	}{
		{"within 24 hours", 12 * time.Hour, 1.0},
		{"within 1 week", 3 * 24 * time.Hour, 0.5},
		{"within 1 month", 15 * 24 * time.Hour, 0.2},
		{"older than 1 month", 60 * 24 * time.Hour, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastUsed := time.Now().Add(-tt.duration)
			score := CalculateTimeScore(lastUsed)

			if score != tt.expected {
				t.Errorf("CalculateTimeScore() = %v, want %v", score, tt.expected)
			}
		})
	}
}

func TestCalculateFrecency(t *testing.T) {
	tests := []struct {
		name     string
		useCount int
		duration time.Duration
		minScore float64
		maxScore float64
	}{
		{"frequently used recently", 10, 1 * time.Hour, 4.0, 5.0},
		{"frequently used long ago", 10, 60 * 24 * time.Hour, 4.0, 4.5},
		{"rarely used recently", 1, 1 * time.Hour, 0.5, 1.0},
		{"rarely used long ago", 1, 60 * 24 * time.Hour, 0.4, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastUsed := time.Now().Add(-tt.duration)
			score := CalculateFrecency(tt.useCount, lastUsed)

			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("CalculateFrecency() = %v, want between %v and %v", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestScoreScripts(t *testing.T) {
	scripts := []NPMScript{
		{Name: "dev", Command: "next dev"},
		{Name: "build", Command: "next build"},
		{Name: "test", Command: "jest"},
	}

	usageStats := []ScriptUsage{
		{
			ScriptName: "dev",
			LastUsed:   time.Now().Add(-1 * time.Hour),
			UseCount:   10,
		},
		{
			ScriptName: "build",
			LastUsed:   time.Now().Add(-24 * time.Hour),
			UseCount:   5,
		},
	}

	scoredScripts := ScoreScripts(scripts, usageStats)

	if len(scoredScripts) != 3 {
		t.Fatalf("expected 3 scored scripts, got %d", len(scoredScripts))
	}

	// Check that scripts are sorted by frecency (descending)
	for i := 0; i < len(scoredScripts)-1; i++ {
		if scoredScripts[i].FrecencyScore < scoredScripts[i+1].FrecencyScore {
			t.Errorf("scripts not sorted by frecency: %v has lower score than %v",
				scoredScripts[i].Script.Name, scoredScripts[i+1].Script.Name)
		}
	}

	// Check that "dev" has highest score (used frequently and recently)
	if scoredScripts[0].Script.Name != "dev" {
		t.Errorf("expected 'dev' to have highest score, got '%s'", scoredScripts[0].Script.Name)
	}

	// Check that "test" has zero score (never used)
	testScript := findScriptByName(scoredScripts, "test")
	if testScript == nil {
		t.Fatal("test script not found")
	}
	if testScript.FrecencyScore != 0.0 {
		t.Errorf("expected test script to have 0 score, got %v", testScript.FrecencyScore)
	}
}

func TestGetMostFrecent(t *testing.T) {
	scripts := []NPMScript{
		{Name: "dev", Command: "next dev"},
		{Name: "build", Command: "next build"},
	}

	usageStats := []ScriptUsage{
		{
			ScriptName: "dev",
			LastUsed:   time.Now().Add(-1 * time.Hour),
			UseCount:   10,
		},
		{
			ScriptName: "build",
			LastUsed:   time.Now().Add(-24 * time.Hour),
			UseCount:   5,
		},
	}

	scoredScripts := ScoreScripts(scripts, usageStats)
	mostFrecent := GetMostFrecent(scoredScripts)

	if mostFrecent == nil {
		t.Fatal("expected most frecent script, got nil")
	}

	if mostFrecent.Script.Name != "dev" {
		t.Errorf("expected 'dev' as most frecent, got '%s'", mostFrecent.Script.Name)
	}
}

func TestGetMostFrecentWithNoHistory(t *testing.T) {
	scripts := []NPMScript{
		{Name: "dev", Command: "next dev"},
		{Name: "build", Command: "next build"},
	}

	usageStats := []ScriptUsage{}

	scoredScripts := ScoreScripts(scripts, usageStats)
	mostFrecent := GetMostFrecent(scoredScripts)

	if mostFrecent != nil {
		t.Error("expected nil for most frecent with no history")
	}
}

func TestGetMostFrecentEmptyScripts(t *testing.T) {
	scoredScripts := []ScoredScript{}
	mostFrecent := GetMostFrecent(scoredScripts)

	if mostFrecent != nil {
		t.Error("expected nil for empty scripts")
	}
}

func TestGetFrecencyStars(t *testing.T) {
	tests := []struct {
		score    float64
		expected string
	}{
		{15.0, "★★★★★"},
		{10.0, "★★★★★"},
		{8.0, "★★★★☆"},
		{5.0, "★★★☆☆"},
		{2.0, "★★☆☆☆"},
		{0.5, "★☆☆☆☆"},
		{0.0, "☆☆☆☆☆"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			stars := GetFrecencyStars(tt.score)
			if stars != tt.expected {
				t.Errorf("GetFrecencyStars(%v) = %s, want %s", tt.score, stars, tt.expected)
			}
		})
	}
}

func TestScoreScriptsWithPartialHistory(t *testing.T) {
	scripts := []NPMScript{
		{Name: "dev", Command: "next dev"},
		{Name: "build", Command: "next build"},
		{Name: "test", Command: "jest"},
		{Name: "lint", Command: "eslint"},
	}

	usageStats := []ScriptUsage{
		{
			ScriptName: "dev",
			LastUsed:   time.Now().Add(-1 * time.Hour),
			UseCount:   20,
		},
		{
			ScriptName: "test",
			LastUsed:   time.Now().Add(-48 * time.Hour),
			UseCount:   3,
		},
	}

	scoredScripts := ScoreScripts(scripts, usageStats)

	if len(scoredScripts) != 4 {
		t.Fatalf("expected 4 scored scripts, got %d", len(scoredScripts))
	}

	// Check that scripts without history have 0 score
	buildScript := findScriptByName(scoredScripts, "build")
	lintScript := findScriptByName(scoredScripts, "lint")

	if buildScript == nil || lintScript == nil {
		t.Fatal("missing scripts")
	}

	if buildScript.FrecencyScore != 0.0 {
		t.Errorf("expected build script to have 0 score, got %v", buildScript.FrecencyScore)
	}

	if lintScript.FrecencyScore != 0.0 {
		t.Errorf("expected lint script to have 0 score, got %v", lintScript.FrecencyScore)
	}
}

// Helper function
func findScriptByName(scripts []ScoredScript, name string) *ScoredScript {
	for _, script := range scripts {
		if script.Script.Name == name {
			return &script
		}
	}
	return nil
}
