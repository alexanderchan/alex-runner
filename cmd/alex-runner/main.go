package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	runner "github.com/alexanderchan/alex-runner/internal"
)

func main() {
	var (
		useLast          bool
		searchTerm       string
		listScripts      bool
		resetDir         bool
		resetAll         bool
		showHelp         bool
		usePackageJSON   bool
		useMakefile      bool
		noCache          bool
	)

	flag.BoolVar(&useLast, "l", false, "Use the most frecent script immediately")
	flag.BoolVar(&useLast, "last", false, "Use the most frecent script immediately")
	flag.StringVar(&searchTerm, "s", "", "Search term for script selection")
	flag.StringVar(&searchTerm, "search", "", "Search term for script selection")
	flag.BoolVar(&listScripts, "list", false, "List all scripts with frecency scores")
	flag.BoolVar(&resetDir, "reset", false, "Clear usage history for current directory")
	flag.BoolVar(&resetAll, "global-reset", false, "Clear all usage history")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&usePackageJSON, "use-package-json", false, "Only show package.json scripts (ignore Makefile)")
	flag.BoolVar(&useMakefile, "use-makefile", false, "Only show Makefile targets (ignore package.json)")
	flag.BoolVar(&noCache, "no-cache", false, "Re-detect package manager instead of using cached value")
	flag.Parse()

	// If no flags provided but positional args exist, use first arg as search term
	if searchTerm == "" && !useLast && !listScripts && !resetDir && !resetAll && len(flag.Args()) > 0 {
		searchTerm = flag.Args()[0]
	}

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: failed to get current directory: %v\n", err)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(cwd)
	if err != nil {
		fmt.Printf("Error: failed to get absolute path: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := runner.InitDatabase()
	if err != nil {
		fmt.Printf("Error: failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Handle reset flags
	if resetAll {
		if err := db.ResetAll(); err != nil {
			fmt.Printf("Error: failed to reset all history: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ All usage history cleared")
		os.Exit(0)
	}

	if resetDir {
		if err := db.ResetDirectory(absPath); err != nil {
			fmt.Printf("Error: failed to reset directory history: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Usage history cleared for %s\n", absPath)
		os.Exit(0)
	}

	// Determine which sources to use
	var scripts []runner.NPMScript

	hasMakefile := runner.MakefileExists(absPath)
	hasPackageJSON := runner.PackageJSONExists(absPath)

	// Detect package manager with caching
	var packageManager string
	if !noCache {
		// Try to get from cache first
		cachedPM, err := db.GetCachedPackageManager(absPath)
		if err != nil {
			fmt.Printf("Warning: failed to get cached package manager: %v\n", err)
		}
		if cachedPM != "" {
			packageManager = cachedPM
		}
	}

	// If no cache or --no-cache flag, detect fresh
	if packageManager == "" {
		packageManager = runner.DetectPackageManager(absPath)
		// Cache the detected package manager
		if err := db.SetCachedPackageManager(absPath, packageManager); err != nil {
			fmt.Printf("Warning: failed to cache package manager: %v\n", err)
		}
	}

	// Determine which sources to load based on flags and what's available
	loadMakefile := hasMakefile && !usePackageJSON
	loadPackageJSON := hasPackageJSON && !useMakefile

	// If both flags are set, show error
	if usePackageJSON && useMakefile {
		fmt.Println("Error: Cannot use both --use-package-json and --use-makefile")
		os.Exit(1)
	}

	// Load Makefile targets if needed
	if loadMakefile {
		makeTargets, err := runner.ReadMakefile(absPath)
		if err != nil {
			fmt.Printf("Error reading Makefile: %v\n", err)
			os.Exit(1)
		}
		for _, target := range makeTargets {
			scripts = append(scripts, runner.NPMScript{
				Name:    target.Name,
				Command: target.Command,
				Source:  "make",
			})
		}
	}

	// Load package.json scripts if needed
	if loadPackageJSON {
		pkg, err := runner.ReadPackageJSON(absPath)
		if err != nil {
			fmt.Printf("Error reading package.json: %v\n", err)
			os.Exit(1)
		}
		pkgScripts := runner.GetScripts(pkg)
		for i := range pkgScripts {
			pkgScripts[i].Source = packageManager
		}
		scripts = append(scripts, pkgScripts...)
	}

	// Error if no scripts found
	if len(scripts) == 0 {
		fmt.Println("Error: No Makefile or package.json found in current directory")
		os.Exit(1)
	}

	// Get usage stats
	usageStats, err := db.GetUsageStats(absPath)
	if err != nil {
		fmt.Printf("Error: failed to get usage stats: %v\n", err)
		os.Exit(1)
	}

	// Score and sort scripts
	scoredScripts := runner.ScoreScripts(scripts, usageStats)

	// Handle list flag
	if listScripts {
		runner.PrintScriptsList(scoredScripts, packageManager)
		os.Exit(0)
	}

	var selectedScript *runner.ScoredScript

	// Handle search term with -l flag: "I'm feeling lucky" with search
	if searchTerm != "" && useLast {
		searchResults := runner.SearchScripts(scoredScripts, searchTerm)
		if len(searchResults) == 0 {
			fmt.Printf("No scripts matching '%s' found\n", searchTerm)
			os.Exit(1)
		}
		// Run first match immediately
		selectedScript = &searchResults[0]
		fmt.Printf("Selected: %s → %s\n", selectedScript.Script.Name, selectedScript.Script.Command)
	} else if searchTerm != "" {
		// Search without -l: show custom selector with editable filter pre-populated with search term
		// Use all scripts (not pre-filtered) so user can edit and see different results
		selected, err := runner.ShowScriptSelectionWithFilter(scoredScripts, searchTerm)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		selectedScript = selected
	} else if useLast {
		// -l without search: use most frecent
		mostFrecent := runner.GetMostFrecent(scoredScripts)
		if mostFrecent == nil {
			fmt.Println("No script usage history found. Please select a script:")
			selected, err := runner.ShowScriptSelection(scoredScripts, "")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			selectedScript = selected
		} else {
			selectedScript = mostFrecent
		}
	} else {
		// Default behavior: show interactive selection
		selected, err := runner.ShowScriptSelection(scoredScripts, "")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		selectedScript = selected
	}

	if selectedScript == nil {
		// User cancelled selection - exit cleanly
		os.Exit(0)
	}

	// Record usage
	if err := db.RecordUsage(absPath, selectedScript.Script.Name); err != nil {
		fmt.Printf("Warning: failed to record usage: %v\n", err)
	}

	// Execute script based on its source
	if selectedScript.Script.Source == "make" {
		fmt.Printf("\n🚀 Running: make %s\n\n", selectedScript.Script.Name)
		if err := executeScript("make", selectedScript.Script.Name, false); err != nil {
			fmt.Printf("Error: script execution failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		// For npm/pnpm/yarn
		fmt.Printf("\n🚀 Running: %s run %s\n\n", selectedScript.Script.Source, selectedScript.Script.Name)
		if err := executeScript(selectedScript.Script.Source, selectedScript.Script.Name, true); err != nil {
			fmt.Printf("Error: script execution failed: %v\n", err)
			os.Exit(1)
		}
	}
}

func executeScript(command string, scriptName string, useRun bool) error {
	var cmd *exec.Cmd
	if useRun {
		cmd = exec.Command(command, "run", scriptName)
	} else {
		cmd = exec.Command(command, scriptName)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func printHelp() {
	help := `alex-runner - Frecency-based npm script runner

USAGE:
    alex-runner [FLAGS] [SEARCH_TERM]

FLAGS:
    -l, --last               Run the most frecent script immediately
    -s, --search <term>      Search for scripts matching term
    --list                   List all scripts with frecency scores
    --use-package-json       Only show package.json scripts (ignore Makefile)
    --use-makefile           Only show Makefile targets (ignore package.json)
    --no-cache               Re-detect package manager (ignore cached detection)
    --reset                  Clear usage history for current directory
    --global-reset           Clear all usage history
    -h, --help               Show this help message

EXAMPLES:
    alex-runner                  # Interactive mode with live filtering
    alex-runner build            # Show selector filtered to "build" matches
    alex-runner -l               # "I'm feeling lucky" - run most frecent immediately
    alex-runner -l build         # Lucky + search - run first "build" match
    alex-runner -s test          # Show selector filtered to "test" matches
    alex-runner --list           # Show all scripts with stats
    alex-runner --use-makefile   # Only show Makefile targets
    alex-runner --reset          # Clear history for current project

BEHAVIOR:
    By default, alex-runner will:
    1. Show scripts from both Makefile and package.json (if both exist)
    2. Show interactive script selection with most frecent at the top
    3. Start typing to filter scripts in real-time
    4. Display script names, commands, and source (make/npm/pnpm/yarn)
    5. Track usage to improve suggestions over time

    Use --use-makefile or --use-package-json to filter to a single source.

The tool stores usage data per directory in ~/.config/alex-runner/
`
	fmt.Println(help)
}
