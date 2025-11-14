package runner

import (
	"strings"
	"testing"
)

func TestGenerateBashCompletion(t *testing.T) {
	completion := GenerateBashCompletion()

	// Check for essential components
	tests := []struct {
		name     string
		contains string
	}{
		{"function definition", "_alex_runner_completion()"},
		{"complete command", "complete -F _alex_runner_completion alex-runner"},
		{"flag -l", "-l"},
		{"flag --last", "--last"},
		{"flag -s", "-s"},
		{"flag --search", "--search"},
		{"flag --list", "--list"},
		{"flag --list-names", "--list-names"},
		{"flag --generate-completion", "--generate-completion"},
		{"flag --use-package-json", "--use-package-json"},
		{"flag --use-makefile", "--use-makefile"},
		{"flag --no-cache", "--no-cache"},
		{"flag --reset", "--reset"},
		{"flag --global-reset", "--global-reset"},
		{"double dash handling", "# Handle -- separator"},
		{"search completion", "alex-runner --list-names"},
		{"shell completion choices", "bash zsh fish"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(completion, tt.contains) {
				t.Errorf("Bash completion missing '%s'", tt.contains)
			}
		})
	}

	// Check that it starts with a comment
	if !strings.HasPrefix(completion, "#") {
		t.Error("Bash completion should start with a comment")
	}
}

func TestGenerateZshCompletion(t *testing.T) {
	completion := GenerateZshCompletion()

	// Check for essential components
	tests := []struct {
		name     string
		contains string
	}{
		{"compdef directive", "#compdef alex-runner"},
		{"function definition", "_alex_runner()"},
		{"helper function", "_alex_runner_scripts()"},
		{"flag -l", "-l"},
		{"flag --last", "--last"},
		{"flag -s", "-s"},
		{"flag --search", "--search"},
		{"flag --list", "--list"},
		{"flag --list-names", "--list-names"},
		{"flag --generate-completion", "--generate-completion"},
		{"script completion", "alex-runner --list-names"},
		{"shell choices", ":(bash zsh fish)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(completion, tt.contains) {
				t.Errorf("Zsh completion missing '%s'", tt.contains)
			}
		})
	}

	// Check that it starts with compdef
	if !strings.HasPrefix(completion, "#compdef") {
		t.Error("Zsh completion should start with #compdef")
	}
}

func TestGenerateFishCompletion(t *testing.T) {
	completion := GenerateFishCompletion()

	// Check for essential components
	tests := []struct {
		name     string
		contains string
	}{
		{"comment header", "# fish completion"},
		{"double dash function", "__alex_runner_after_double_dash"},
		{"scripts function", "__alex_runner_scripts"},
		{"complete command", "complete -c alex-runner"},
		{"flag -h", "-s h"},
		{"flag -l", "-s l"},
		{"flag -s", "-s s"},
		{"flag --help", "-l help"},
		{"flag --last", "-l last"},
		{"flag --search", "-l search"},
		{"flag --list", "-l list"},
		{"flag --list-names", "-l list-names"},
		{"flag --generate-completion", "-l generate-completion"},
		{"script completion", "alex-runner --list-names"},
		{"shell choices", "'bash zsh fish'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(completion, tt.contains) {
				t.Errorf("Fish completion missing '%s'", tt.contains)
			}
		})
	}

	// Check that it starts with a comment
	if !strings.HasPrefix(completion, "#") {
		t.Error("Fish completion should start with a comment")
	}
}

func TestGetCompletionInstallInstructions(t *testing.T) {
	tests := []struct {
		shell    string
		contains []string
	}{
		{
			shell: "bash",
			contains: []string{
				"Bash Completion Installation",
				"--generate-completion bash",
				".bashrc",
				".alex-runner-completion.bash",
			},
		},
		{
			shell: "zsh",
			contains: []string{
				"Zsh Completion Installation",
				"--generate-completion zsh",
				".zshrc",
				".alex-runner-completion.zsh",
				"compinit",
			},
		},
		{
			shell: "fish",
			contains: []string{
				"Fish Completion Installation",
				"--generate-completion fish",
				"~/.config/fish/completions",
				"alex-runner.fish",
			},
		},
		{
			shell: "unknown",
			contains: []string{
				"Unknown shell",
				"Supported shells: bash, zsh, fish",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			instructions := GetCompletionInstallInstructions(tt.shell)
			for _, expected := range tt.contains {
				if !strings.Contains(instructions, expected) {
					t.Errorf("Instructions for %s missing '%s'", tt.shell, expected)
				}
			}
		})
	}
}

func TestListCompletions(t *testing.T) {
	list := ListCompletions()

	expectedStrings := []string{
		"bash",
		"zsh",
		"fish",
		"--generate-completion",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(list, expected) {
			t.Errorf("Completion list missing '%s'", expected)
		}
	}
}

// Test that all three completion formats are non-empty and reasonably sized
func TestCompletionFormatsLength(t *testing.T) {
	tests := []struct {
		name       string
		completion string
		minLength  int
	}{
		{"bash", GenerateBashCompletion(), 500},
		{"zsh", GenerateZshCompletion(), 400},
		{"fish", GenerateFishCompletion(), 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.completion) < tt.minLength {
				t.Errorf("%s completion too short: got %d bytes, want at least %d",
					tt.name, len(tt.completion), tt.minLength)
			}
		})
	}
}

// Test that completions don't contain any obvious syntax errors
func TestCompletionSyntax(t *testing.T) {
	tests := []struct {
		name       string
		completion string
		mustNotContain []string
	}{
		{
			name:       "bash no unmatched quotes",
			completion: GenerateBashCompletion(),
			mustNotContain: []string{},
		},
		{
			name:       "zsh no unmatched quotes",
			completion: GenerateZshCompletion(),
			mustNotContain: []string{},
		},
		{
			name:       "fish no unmatched quotes",
			completion: GenerateFishCompletion(),
			mustNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check for balanced single quotes (excluding comments)
			lines := strings.Split(tt.completion, "\n")
			for i, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "#") {
					continue // Skip comments
				}
				singleQuotes := strings.Count(line, "'")
				if singleQuotes%2 != 0 {
					t.Errorf("Line %d has unbalanced single quotes: %s", i+1, line)
				}
			}
		})
	}
}

// Test case insensitivity of shell parameter
func TestGetCompletionInstallInstructionsCaseInsensitive(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"bash", "Bash"},
		{"BASH", "Bash"},
		{"Bash", "Bash"},
		{"zsh", "Zsh"},
		{"ZSH", "Zsh"},
		{"fish", "Fish"},
		{"FISH", "Fish"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			instructions := GetCompletionInstallInstructions(tt.input)
			if !strings.Contains(instructions, tt.expected) {
				t.Errorf("Case-insensitive test failed: input=%s, expected to contain '%s'",
					tt.input, tt.expected)
			}
		})
	}
}
