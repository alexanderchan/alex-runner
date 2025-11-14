package runner

import (
	"reflect"
	"testing"
)

func TestBuildScriptArgs_NPM_WithArgs(t *testing.T) {
	result := BuildScriptArgs(BuildScriptArgsParams{
		Command:        "npm",
		ScriptName:     "test",
		UseRun:         true,
		AdditionalArgs: []string{"--testPathPattern", "some/path", "--verbose"},
	})

	expected := []string{"run", "test", "--", "--testPathPattern", "some/path", "--verbose"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBuildScriptArgs_NPM_NoArgs(t *testing.T) {
	result := BuildScriptArgs(BuildScriptArgsParams{
		Command:        "npm",
		ScriptName:     "test",
		UseRun:         true,
		AdditionalArgs: []string{},
	})

	expected := []string{"run", "test"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBuildScriptArgs_PNPM_WithArgs(t *testing.T) {
	result := BuildScriptArgs(BuildScriptArgsParams{
		Command:        "pnpm",
		ScriptName:     "test",
		UseRun:         true,
		AdditionalArgs: []string{"--reporter", "verbose"},
	})

	// pnpm should NOT have '--' separator
	expected := []string{"run", "test", "--reporter", "verbose"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBuildScriptArgs_PNPM_NoArgs(t *testing.T) {
	result := BuildScriptArgs(BuildScriptArgsParams{
		Command:        "pnpm",
		ScriptName:     "build",
		UseRun:         true,
		AdditionalArgs: []string{},
	})

	expected := []string{"run", "build"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBuildScriptArgs_Yarn_WithArgs(t *testing.T) {
	result := BuildScriptArgs(BuildScriptArgsParams{
		Command:        "yarn",
		ScriptName:     "test",
		UseRun:         true,
		AdditionalArgs: []string{"--watch"},
	})

	// yarn should NOT have '--' separator
	expected := []string{"run", "test", "--watch"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBuildScriptArgs_Make_WithArgs(t *testing.T) {
	result := BuildScriptArgs(BuildScriptArgsParams{
		Command:        "make",
		ScriptName:     "test",
		UseRun:         false,
		AdditionalArgs: []string{"-j", "4"},
	})

	expected := []string{"test", "-j", "4"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBuildScriptArgs_Make_NoArgs(t *testing.T) {
	result := BuildScriptArgs(BuildScriptArgsParams{
		Command:        "make",
		ScriptName:     "build",
		UseRun:         false,
		AdditionalArgs: []string{},
	})

	expected := []string{"build"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseArgs_WithSeparator(t *testing.T) {
	args := []string{"-l", "test", "--", "--testPathPattern", "some/path"}
	before, after := ParseArgs(args)

	expectedBefore := []string{"-l", "test"}
	expectedAfter := []string{"--testPathPattern", "some/path"}

	if !reflect.DeepEqual(before, expectedBefore) {
		t.Errorf("Expected before %v, got %v", expectedBefore, before)
	}
	if !reflect.DeepEqual(after, expectedAfter) {
		t.Errorf("Expected after %v, got %v", expectedAfter, after)
	}
}

func TestParseArgs_NoSeparator(t *testing.T) {
	args := []string{"-l", "test", "--verbose"}
	before, after := ParseArgs(args)

	if !reflect.DeepEqual(before, args) {
		t.Errorf("Expected before %v, got %v", args, before)
	}
	if after != nil {
		t.Errorf("Expected after to be nil, got %v", after)
	}
}

func TestParseArgs_SeparatorAtStart(t *testing.T) {
	args := []string{"--", "--watch", "--verbose"}
	before, after := ParseArgs(args)

	expectedBefore := []string{}
	expectedAfter := []string{"--watch", "--verbose"}

	if !reflect.DeepEqual(before, expectedBefore) {
		t.Errorf("Expected before %v, got %v", expectedBefore, before)
	}
	if !reflect.DeepEqual(after, expectedAfter) {
		t.Errorf("Expected after %v, got %v", expectedAfter, after)
	}
}

func TestParseArgs_SeparatorAtEnd(t *testing.T) {
	args := []string{"-l", "test", "--"}
	before, after := ParseArgs(args)

	expectedBefore := []string{"-l", "test"}
	expectedAfter := []string{}

	if !reflect.DeepEqual(before, expectedBefore) {
		t.Errorf("Expected before %v, got %v", expectedBefore, before)
	}
	if !reflect.DeepEqual(after, expectedAfter) {
		t.Errorf("Expected after %v, got %v", expectedAfter, after)
	}
}

func TestParseArgs_EmptyInput(t *testing.T) {
	args := []string{}
	before, after := ParseArgs(args)

	if !reflect.DeepEqual(before, args) {
		t.Errorf("Expected before %v, got %v", args, before)
	}
	if after != nil {
		t.Errorf("Expected after to be nil, got %v", after)
	}
}

func TestBuildScriptArgs_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		params   BuildScriptArgsParams
		expected []string
	}{
		{
			name: "npm with multiple args",
			params: BuildScriptArgsParams{
				Command:        "npm",
				ScriptName:     "test",
				UseRun:         true,
				AdditionalArgs: []string{"--coverage", "--watch"},
			},
			expected: []string{"run", "test", "--", "--coverage", "--watch"},
		},
		{
			name: "pnpm with single arg",
			params: BuildScriptArgsParams{
				Command:        "pnpm",
				ScriptName:     "dev",
				UseRun:         true,
				AdditionalArgs: []string{"--port=3000"},
			},
			expected: []string{"run", "dev", "--port=3000"},
		},
		{
			name: "yarn with no args",
			params: BuildScriptArgsParams{
				Command:        "yarn",
				ScriptName:     "build",
				UseRun:         true,
				AdditionalArgs: nil,
			},
			expected: []string{"run", "build"},
		},
		{
			name: "make with args containing equals",
			params: BuildScriptArgsParams{
				Command:        "make",
				ScriptName:     "install",
				UseRun:         false,
				AdditionalArgs: []string{"PREFIX=/usr/local"},
			},
			expected: []string{"install", "PREFIX=/usr/local"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildScriptArgs(tt.params)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
