package runner

import (
	"testing"
	"time"
)

// TestMultiWordThreshold demonstrates the threshold filtering behavior
func TestMultiWordThreshold(t *testing.T) {
	now := time.Now()
	scripts := []ScoredScript{
		{
			Script:        NPMScript{Name: "hello", Command: "echo '👋 Hello from Makefile!'", Source: "make"},
			FrecencyScore: 8.0,
			LastUsed:      &now,
			UseCount:      8,
		},
		{
			Script:        NPMScript{Name: "install", Command: "go install ./cmd/alex-runner", Source: "make"},
			FrecencyScore: 11.0,
			LastUsed:      &now,
			UseCount:      11,
		},
		{
			Script:        NPMScript{Name: "test", Command: "go test ./internal/... -v", Source: "make"},
			FrecencyScore: 1.0,
			LastUsed:      &now,
			UseCount:      1,
		},
		{
			Script:        NPMScript{Name: "build", Command: "go build -o alex-runner ./cmd/alex-runner", Source: "make"},
			FrecencyScore: 10.0,
			LastUsed:      &now,
			UseCount:      10,
		},
		{
			Script:        NPMScript{Name: "docker-build", Command: "docker build -t myapp .", Source: "make"},
			FrecencyScore: 5.0,
			LastUsed:      &now,
			UseCount:      5,
		},
		{
			Script:        NPMScript{Name: "hello-docker", Command: "docker run hello-world", Source: "make"},
			FrecencyScore: 3.0,
			LastUsed:      &now,
			UseCount:      3,
		},
	}

	t.Run("hello docker should filter out weak matches", func(t *testing.T) {
		results := SearchScripts(scripts, "hello docker")

		// Should only return scripts that strongly match BOTH terms
		t.Logf("Found %d results for 'hello docker'", len(results))
		for i, result := range results {
			t.Logf("  [%d] %s (command: %s)", i, result.Script.Name, result.Script.Command)
		}

		// Should find hello-docker (contains both terms)
		foundHelloDocker := false
		foundPlainHello := false

		for _, result := range results {
			if result.Script.Name == "hello-docker" {
				foundHelloDocker = true
			}
			if result.Script.Name == "hello" {
				foundPlainHello = true
			}
		}

		if !foundHelloDocker {
			t.Error("Expected to find 'hello-docker' in results")
		}

		if foundPlainHello {
			t.Error("Should NOT find plain 'hello' in results (doesn't contain 'docker')")
		}

		// With the threshold, we should have fewer results than total scripts
		if len(results) >= len(scripts) {
			t.Errorf("Threshold should filter out weak matches. Got %d results out of %d scripts", len(results), len(scripts))
		}
	})

	t.Run("single word hello should still work", func(t *testing.T) {
		results := SearchScripts(scripts, "hello")

		foundHello := false
		for _, result := range results {
			if result.Script.Name == "hello" {
				foundHello = true
				break
			}
		}

		if !foundHello {
			t.Error("Single word search 'hello' should find 'hello' script")
		}
	})
}
