package runner

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MakeTarget represents a Makefile target
type MakeTarget struct {
	Name    string
	Command string
}

// ReadMakefile reads and parses targets from a Makefile
func ReadMakefile(directory string) ([]MakeTarget, error) {
	makefilePath := filepath.Join(directory, "Makefile")

	file, err := os.Open(makefilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var targets []MakeTarget
	scanner := bufio.NewScanner(file)

	// Regex to match target definitions: "targetname:" or "targetname: dependencies"
	targetRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+):\s*(.*)$`)

	var currentTarget *MakeTarget

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Check if this is a target definition
		if matches := targetRegex.FindStringSubmatch(line); matches != nil {
			targetName := matches[1]

			// Skip special targets like .PHONY
			if strings.HasPrefix(targetName, ".") {
				continue
			}

			// Save previous target if exists
			if currentTarget != nil && currentTarget.Command != "" {
				targets = append(targets, *currentTarget)
			}

			// Start new target
			currentTarget = &MakeTarget{
				Name:    targetName,
				Command: "",
			}
		} else if currentTarget != nil && strings.HasPrefix(line, "\t") {
			// This is a command line (starts with tab)
			command := strings.TrimPrefix(line, "\t")
			// Remove @ prefix if present (suppresses echo)
			command = strings.TrimPrefix(command, "@")

			if currentTarget.Command != "" {
				currentTarget.Command += " && " + command
			} else {
				currentTarget.Command = command
			}
		}
	}

	// Add last target
	if currentTarget != nil && currentTarget.Command != "" {
		targets = append(targets, *currentTarget)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}

// MakefileExists checks if a Makefile exists in the directory
func MakefileExists(directory string) bool {
	makefilePath := filepath.Join(directory, "Makefile")
	_, err := os.Stat(makefilePath)
	return err == nil
}
