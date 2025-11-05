package runner

import (
	"sort"
	"time"
)

const (
	frequencyWeight = 0.4
	recencyWeight   = 0.6
)

type ScoredScript struct {
	Script       NPMScript
	FrecencyScore float64
	LastUsed     *time.Time
	UseCount     int
}

func CalculateTimeScore(lastUsed time.Time) float64 {
	duration := time.Since(lastUsed)

	switch {
	case duration < 24*time.Hour:
		return 1.0
	case duration < 7*24*time.Hour:
		return 0.5
	case duration < 30*24*time.Hour:
		return 0.2
	default:
		return 0.1
	}
}

func CalculateFrecency(useCount int, lastUsed time.Time) float64 {
	timeScore := CalculateTimeScore(lastUsed)
	frequencyScore := float64(useCount)

	return (frequencyScore * frequencyWeight) + (timeScore * recencyWeight)
}

func ScoreScripts(scripts []NPMScript, usageStats []ScriptUsage) []ScoredScript {
	// Create a map of script names to usage stats for quick lookup
	usageMap := make(map[string]ScriptUsage)
	for _, usage := range usageStats {
		usageMap[usage.ScriptName] = usage
	}

	scoredScripts := make([]ScoredScript, 0, len(scripts))

	for _, script := range scripts {
		scored := ScoredScript{
			Script: script,
		}

		if usage, exists := usageMap[script.Name]; exists {
			scored.FrecencyScore = CalculateFrecency(usage.UseCount, usage.LastUsed)
			scored.LastUsed = &usage.LastUsed
			scored.UseCount = usage.UseCount
		} else {
			// New script with no history
			scored.FrecencyScore = 0.0
			scored.LastUsed = nil
			scored.UseCount = 0
		}

		scoredScripts = append(scoredScripts, scored)
	}

	// Sort by frecency score descending (highest score first)
	sort.Slice(scoredScripts, func(i, j int) bool {
		return scoredScripts[i].FrecencyScore > scoredScripts[j].FrecencyScore
	})

	return scoredScripts
}

func GetMostFrecent(scoredScripts []ScoredScript) *ScoredScript {
	if len(scoredScripts) == 0 {
		return nil
	}

	// Scripts are already sorted by frecency descending
	// Return the first one (highest score)
	if scoredScripts[0].FrecencyScore > 0 {
		return &scoredScripts[0]
	}

	// If no scripts have been used, return nil
	return nil
}

func GetFrecencyStars(score float64) string {
	// Convert score to 0-5 stars
	// Typical scores range from 0 to ~10+
	// We'll map: 0-1 = 1 star, 1-3 = 2 stars, 3-6 = 3 stars, 6-10 = 4 stars, 10+ = 5 stars
	switch {
	case score >= 10:
		return "★★★★★"
	case score >= 6:
		return "★★★★☆"
	case score >= 3:
		return "★★★☆☆"
	case score >= 1:
		return "★★☆☆☆"
	case score > 0:
		return "★☆☆☆☆"
	default:
		return "☆☆☆☆☆"
	}
}
