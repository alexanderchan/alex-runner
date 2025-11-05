package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PackageJSON struct {
	Name    string            `json:"name"`
	Scripts map[string]string `json:"scripts"`
}

type NPMScript struct {
	Name    string
	Command string
	Source  string // "make", "npm", "yarn", "pnpm", etc.
}

// GetGitRoot returns the root of the git repository, or the current directory if not in a git repo
func GetGitRoot(directory string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = directory
	output, err := cmd.Output()
	if err != nil {
		// Not in a git repo, return the current directory
		return directory
	}
	return strings.TrimSpace(string(output))
}

// DetectPackageManager checks for lock files starting from git root to determine package manager
// This is more robust for mono-repos and sub-packages
func DetectPackageManager(directory string) string {
	gitRoot := GetGitRoot(directory)

	// Check git root first for lock files (handles mono-repos)
	if _, err := os.Stat(filepath.Join(gitRoot, "yarn.lock")); err == nil {
		return "yarn"
	}
	if _, err := os.Stat(filepath.Join(gitRoot, "pnpm-lock.yaml")); err == nil {
		return "pnpm"
	}
	if _, err := os.Stat(filepath.Join(gitRoot, "package-lock.json")); err == nil {
		return "npm"
	}

	// If no lock file in git root, check current directory
	if _, err := os.Stat(filepath.Join(directory, "yarn.lock")); err == nil {
		return "yarn"
	}
	if _, err := os.Stat(filepath.Join(directory, "pnpm-lock.yaml")); err == nil {
		return "pnpm"
	}
	if _, err := os.Stat(filepath.Join(directory, "package-lock.json")); err == nil {
		return "npm"
	}

	// If package.json exists but no lock file, default to pnpm (as per user's shell script)
	if PackageJSONExists(directory) || PackageJSONExists(gitRoot) {
		return "pnpm"
	}

	// Final fallback
	return "npm"
}

func PackageJSONExists(directory string) bool {
	packageJSONPath := filepath.Join(directory, "package.json")
	_, err := os.Stat(packageJSONPath)
	return err == nil
}

func ReadPackageJSON(directory string) (*PackageJSON, error) {
	packageJSONPath := filepath.Join(directory, "package.json")

	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no package.json found in current directory")
		}
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	if len(pkg.Scripts) == 0 {
		return nil, fmt.Errorf("no scripts found in package.json")
	}

	return &pkg, nil
}

func GetScripts(pkg *PackageJSON) []NPMScript {
	scripts := make([]NPMScript, 0, len(pkg.Scripts))
	for name, command := range pkg.Scripts {
		scripts = append(scripts, NPMScript{
			Name:    name,
			Command: command,
			Source:  "", // Will be set by caller
		})
	}
	return scripts
}
