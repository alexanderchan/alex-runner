# alex-runner

A smart, frecency-based npm script runner that learns which scripts you use most often and suggests them automatically.

![Alex Runner Logo](public/images/logo.jpg)

## Problem

Sometimes there's a lot of scripts but you only really use a few of them and forget their names. alex-runner helps you remember which scripts you use most often and suggests them automatically.

## Features

- **Script Pinning**: Pin your most important scripts to always appear first (📌)
  - Use `--pin <script>` from CLI or press `alt-p` in the UI
  - Handles duplicate script names across Makefile and package.json
- **Frecency-based suggestions**: Combines frequency and recency to suggest the scripts you're most likely to need
- **Live filtering**: Type to instantly filter scripts - no special keys needed
- **Beautiful TUI**: Powered by Bubble Tea with syntax highlighting and clear command previews
- **Shell completion**: Tab completion for bash/zsh/fish with frecency-aware script suggestions
- **Multi-package manager**: Automatically detects npm, pnpm, or yarn
- **Makefile support**: Run Makefile targets alongside npm scripts
- **Per-directory tracking**: Each project has its own usage history
- **Source-aware tracking**: Scripts from Makefile and package.json are tracked separately
- **Fuzzy search**: Quickly find scripts by name or command content
- **Smart search ranking**: 6-tier priority system from exact matches to fuzzy command matches
- **Zero configuration**: Just install and run

## Demo Flow

![Alex Runner Demo](public/images/demo.gif)

When you run `alex-runner`:

1. **Interactive selector**: Shows all scripts with most frecent at the top
   ```
   📦 Select an npm script to run (type to filter)

   > dev
       → next dev [★★★★★ 24 runs, 2h ago]

     build
       → next build [★★☆☆☆ 5 runs, 1d ago]

     typecheck
       → tsc --noEmit [★☆☆☆☆ 2 runs, 3h ago]
   ```

2. **Type to filter**: Start typing (e.g., "bui") to instantly filter the list
3. **Arrow keys**: Navigate up/down through options
4. **Enter**: Run the selected script

## Installation

### Via Go Install (Recommended)

```bash
go install github.com/alexanderchan/alex-runner/cmd/alex-runner@latest
```

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Download Pre-built Binary

Download the latest release for your platform from the [Releases page](https://github.com/alexanderchan/alex-runner/releases).

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/alexanderchan/alex-runner/releases/latest/download/alex-runner_Darwin_arm64.tar.gz | tar xz
sudo mv alex-runner /usr/local/bin/
# Remove quarantine attribute (required for unsigned binaries)
xattr -d com.apple.quarantine /usr/local/bin/alex-runner
```

**Linux:**
```bash
# For ARM64
curl -L https://github.com/alexanderchan/alex-runner/releases/latest/download/alex-runner_Linux_arm64.tar.gz | tar xz
sudo mv alex-runner /usr/local/bin/

# For x86_64
curl -L https://github.com/alexanderchan/alex-runner/releases/latest/download/alex-runner_Linux_x86_64.tar.gz | tar xz
sudo mv alex-runner /usr/local/bin/
```

**Windows:**
Download the `.zip` file from the [Releases page](https://github.com/alexanderchan/alex-runner/releases), extract it, and add the directory to your PATH.

### Manual Build

```bash
git clone https://github.com/alexanderchan/alex-runner
cd alex-runner

# Using npm/pnpm
npm run build
# or
pnpm build

# Or using go directly
go build -o alex-runner ./cmd/alex-runner

# Move to PATH
sudo mv alex-runner /usr/local/bin/
```

### Shell Completion

alex-runner supports tab completion for bash, zsh, and fish shells. Completions are **frecency-aware**, meaning your most-used scripts appear first!

#### Bash

```bash
# Generate and save completion script
alex-runner --generate-completion bash > ~/.alex-runner-completion.bash

# Add to your ~/.bashrc
echo 'source ~/.alex-runner-completion.bash' >> ~/.bashrc

# Reload your shell
source ~/.bashrc
```

**System-wide installation:**
```bash
sudo alex-runner --generate-completion bash > /etc/bash_completion.d/alex-runner
```

#### Zsh

```bash
# Generate and save completion script
alex-runner --generate-completion zsh > ~/.alex-runner-completion.zsh

# Add to your ~/.zshrc (MUST be AFTER compinit)
echo 'source ~/.alex-runner-completion.zsh' >> ~/.zshrc

# Reload your shell
source ~/.zshrc
```

**Important:** The completion must be sourced **after** `compinit` is called in your `.zshrc`. Example:

```zsh
# In your ~/.zshrc:
autoload -Uz compinit && compinit

# AFTER compinit, source alex-runner completion
source ~/.alex-runner-completion.zsh
```

**Using a completion directory:**
```bash
# Create directory if it doesn't exist
mkdir -p ~/.zsh/completions

# Generate completion file
alex-runner --generate-completion zsh > ~/.zsh/completions/_alex-runner

# Add to ~/.zshrc (before compinit)
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -Uz compinit && compinit' >> ~/.zshrc

# Reload shell
exec zsh
```

#### Fish

```bash
# Create completions directory if needed
mkdir -p ~/.config/fish/completions

# Generate completion script
alex-runner --generate-completion fish > ~/.config/fish/completions/alex-runner.fish

# Reload completions (or restart fish)
fish_update_completions
```

#### What Gets Completed

- **Flags**: `-l`, `--last`, `-s`, `--search`, `--list`, `--reset`, etc.
- **Script names**: Dynamically fetched from current directory (sorted by frecency!)
- **Shell types**: After `--generate-completion`, suggests `bash`, `zsh`, `fish`
- **Smart context**: After `--`, completions stop (those are script arguments)

#### Testing Completions

```bash
# Tab completion for flags
alex-runner --<TAB>
# Shows: --last, --list, --search, --generate-completion, etc.

# Tab completion for scripts (frecency-aware!)
alex-runner <TAB>
# Shows: dev, build, test, etc. (your most-used first)

# After --search flag
alex-runner --search <TAB>
# Shows: dev, build, test, etc.
```

**Note:** Completions work in standard terminals (iTerm2, Terminal.app, Alacritty, etc.). Some terminals with custom completion systems (like Warp) may not support standard shell completions.

#### Completion for Aliases

If you have aliases for alex-runner (e.g., `alias rr="alex-runner"`), you need to register completion for them:

**Zsh:**
```zsh
# In your ~/.zshrc, after sourcing the completion:
alias rr="alex-runner"
alias rrl="alex-runner -l"

# Register completions for your aliases
compdef _alex_runner rr
compdef _alex_runner rrl
```

**Bash:**
```bash
# In your ~/.bashrc:
alias rr="alex-runner"
alias rrl="alex-runner -l"

# Register completions for your aliases
complete -F _alex_runner_completion rr
complete -F _alex_runner_completion rrl
```

**Fish:**
```fish
# In your ~/.config/fish/config.fish:
alias rr="alex-runner"
alias rrl="alex-runner -l"

# Fish automatically handles completion for aliases, no extra setup needed!
```

## Upgrading

### Upgrading to v0.2.0

Version 0.2.0 includes a breaking change to the database schema to support script pinning and better source tracking.

**Migration Required**: If upgrading from v0.1.0, you need to clear your database:

```bash
# Remove the old database
rm ~/.config/alex-runner/alex-runner.sqlite.db

# alex-runner will automatically create a new database with the updated schema on next run
```

**What changes:**
- Your usage history will be reset
- Scripts from Makefile and package.json are now tracked separately
- Each script+source combination has its own frecency score and pin status

**What's new:**
- 📌 Pin scripts to always appear first
- Better handling of duplicate script names from different sources
- More accurate frecency tracking per source

If you have critical usage history you want to preserve, stay on v0.1.0. Otherwise, the fresh start with improved tracking is recommended.

## Usage

### Interactive Mode (Default)

```bash
cd your-project
alex-runner
```

This will:
1. Show interactive selector with most frecent script at the top
2. Type to instantly filter the list (no special keys needed)
3. Use arrow keys to navigate
4. Display both script names and their actual commands
5. Track your selection for future use

### "I'm Feeling Lucky" Mode

```bash
alex-runner -l
# or
alex-runner --last
```

Immediately runs the most frecent script without prompting - perfect for when you know you want to run the same thing again!

### Search for Scripts

```bash
# Quick search (positional argument)
alex-runner build

# Or using flag
alex-runner -s build
alex-runner --search build
```

Shows an interactive selector filtered to scripts matching "build". You confirm your choice before running - safe and fast!

### Passing Arguments to Scripts

You can pass additional arguments to scripts using the `--` separator:

```bash
# Pass arguments in interactive mode
alex-runner -- --watch --verbose

# Pass arguments with "feeling lucky" mode
alex-runner -l test -- --reporter=verbose

# Pass arguments with search
alex-runner test -- --testPathPattern some/path --coverage
```

**How it works:**
- **npm**: Automatically adds `--` separator (required by npm)
  ```bash
  # alex-runner runs: npm run test -- --watch
  ```
- **pnpm/yarn**: Passes arguments directly (no `--` needed)
  ```bash
  # alex-runner runs: pnpm run test --watch
  ```
- **make**: Appends arguments directly to the target
  ```bash
  # alex-runner runs: make test --verbose
  ```

This is especially useful for test runners, dev servers, and build tools that accept configuration flags.

### List All Scripts

```bash
alex-runner --list
```

Displays all scripts with their frecency scores and usage stats.

### Using Makefile Targets

alex-runner automatically detects and includes Makefile targets:

```bash
# Shows both package.json scripts AND Makefile targets
alex-runner

# Show only Makefile targets
alex-runner --use-makefile

# Show only package.json scripts
alex-runner --use-package-json
```

Makefile targets are displayed with a "make" indicator and run with `make target-name` instead of the package manager.

### Pin Scripts

Pin your most important scripts to always appear first, regardless of frecency:

```bash
# Pin a script
alex-runner --pin dev

# If the script exists in both Makefile and package.json, you'll be prompted:
Multiple scripts found with name 'dev':
  1. dev (make)
  2. dev (pnpm)
Select which one to pin (1-2, or 'all'): 1

# Unpin a script
alex-runner --unpin dev

# Toggle pin in interactive mode
# Press alt-p (or option-p on Mac) to toggle pin for the selected script
```

Pinned scripts:
- Show a 📌 indicator
- Always appear at the top of the list
- Are sorted by frecency among themselves
- Are tracked per source (Makefile vs package.json)

### Reset History

```bash
# Clear history for current directory
alex-runner --reset

# Clear all history
alex-runner --global-reset
```

## Frecency Algorithm

alex-runner uses a frecency algorithm to rank scripts based on both frequency and recency:

```
frecency_score = (use_count × 0.4) + (time_score × 0.6)
```

Time scores:
- Last 24 hours: 1.0
- Last week: 0.5
- Last month: 0.2
- Older: 0.1

This means recently used scripts get a boost, but frequently used scripts remain relevant.

## Package Manager Detection

alex-runner automatically detects your package manager by searching for lock files (checks git root first, then current directory):

1. `yarn.lock` found → uses `yarn`
2. `pnpm-lock.yaml` found → uses `pnpm`
3. `package-lock.json` found → uses `npm`
4. `package.json` only (no lock file) → defaults to `pnpm`
5. No files found → falls back to `npm`

Detection results are cached per directory for performance. Use `--no-cache` to force re-detection.

## Data Storage

Usage data is stored in:
```
~/.config/alex-runner/alex-runner.sqlite.db
```

Each directory's script usage is tracked separately, so you get project-specific suggestions.

## UI Colors

- **Script name**: White + Bold
- **Command preview** (`→ ...`): Light gray (readable)
- **Metadata** (stars, run count, time): Dark gray (subtle)
- **Selected item**: Highlighted

## Examples

### First time in a project

```bash
$ alex-runner
📦 Select an npm script to run

  start
    → next start [never used]

  dev
    → next dev [never used]

  build
    → next build [never used]
```

### After using "dev" several times

```bash
$ alex-runner
📦 Select an npm script to run (type to filter)

> dev
    → next dev [★★★★★ 24 runs, 2h ago]

  build
    → next build [★★☆☆☆ 5 runs, 1d ago]

  typecheck
    → tsc --noEmit [★☆☆☆☆ 2 runs, 3h ago]
```

Most frecent at top, type to filter, arrow keys to navigate.

### Searching for scripts

```bash
$ alex-runner type
📦 Select an npm script to run (type to filter)

> typecheck
    → tsc --noEmit [★☆☆☆☆ 2 runs, 3h ago]

[Shows only matching scripts, you press Enter to confirm]
```

Quick and simple - filters to matching scripts so you can confirm before running!

### Listing all scripts

```bash
$ alex-runner --list

Available npm scripts (sorted by frecency):

dev
  → next dev [★★★★★ 24 runs, 2h ago]
  Run with: pnpm run dev

typecheck
  → tsc --noEmit [★★★☆☆ 12 runs, 3h ago]
  Run with: pnpm run typecheck

build
  → next build [★★☆☆☆ 5 runs, 1d ago]
  Run with: pnpm run build
```

## Command-Line Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--last` | `-l` | boolean | false | "I'm feeling lucky" - run most frecent immediately |
| `--search` | `-s` | string | "" | Show selector filtered to search term |
| `--list` | | boolean | false | List all scripts with frecency scores |
| `--list-names` | | boolean | false | List script names only (used for shell completion) |
| `--generate-completion` | | string | "" | Generate shell completion script (bash\|zsh\|fish) |
| `--reset` | | boolean | false | Clear usage history for current directory |
| `--global-reset` | | boolean | false | Clear all usage history |
| `--use-package-json` | | boolean | false | Only show package.json scripts (ignore Makefile) |
| `--use-makefile` | | boolean | false | Only show Makefile targets (ignore package.json) |
| `--no-cache` | | boolean | false | Re-detect package manager instead of using cached detection |
| `--help` | `-h` | boolean | false | Show help message |
| (positional arg) | | string | "" | Same as `--search` - `alex-runner build` |
| `--` | | separator | - | Pass additional arguments to the script (e.g., `alex-runner test -- --watch`) |

## Configuration & Advanced Options

### Frecency Algorithm Parameters

The ranking algorithm uses these weights (defined in source):

```go
frecency_score = (use_count × 0.4) + (time_score × 0.6)
```

**Time-based scores:**
| Duration | Score | Impact |
|----------|-------|--------|
| Last 24 hours | 1.0 | Maximum recency boost |
| Last week (7 days) | 0.5 | Medium recency boost |
| Last month (30 days) | 0.2 | Low recency boost |
| Older than 30 days | 0.1 | Minimal recency boost |

**Star ratings:**
| Frecency Score | Stars | Visual |
|---------------|-------|--------|
| ≥ 10 | 5 stars | ★★★★★ |
| ≥ 6 | 4.5 stars | ★★★★☆ |
| ≥ 3 | 3 stars | ★★★☆☆ |
| ≥ 1 | 2 stars | ★★☆☆☆ |
| > 0 | 1 star | ★☆☆☆☆ |
| = 0 | 0 stars | ☆☆☆☆☆ (never used) |

### Search Ranking System

When you filter scripts, they're ranked by match quality:

| Priority | Match Type | Rank | Example |
|----------|-----------|------|---------|
| 1 | Exact name match | 1000 | "build" → "build" |
| 2 | Name prefix | 500 | "bui" → "build" |
| 3 | Name substring | 300 | "ild" in "build" |
| 4 | Fuzzy name match | 200 | "bld" → "build" |
| 5 | Command substring | 100 | "tsc" in "tsc --noEmit" |
| 6 | Fuzzy command match | 50 | Fuzzy match in command text |

Scripts with the same rank are then sorted by frecency score.

### Package Manager Detection

Detection happens in this order (searches git root first, then current directory):

1. **yarn.lock** found → uses `yarn`
2. **pnpm-lock.yaml** found → uses `pnpm`
3. **package-lock.json** found → uses `npm`
4. **package.json** only (no lock file) → defaults to `pnpm`
5. **No files found** → falls back to `npm`

**Cache behavior**: Detection result is cached per directory. Use `--no-cache` to force re-detection.

### Makefile Support

alex-runner can also run Makefile targets alongside npm scripts:

**Parsing rules:**
- Targets must match pattern: `targetname: [dependencies]`
- Commands must be indented with TAB characters
- Comments (`#`) and `.PHONY` targets are ignored
- The `@` prefix (echo suppression) is automatically removed
- Multiple commands for one target are combined with `&&`

**Filtering:**
- By default, shows both package.json scripts AND Makefile targets
- Use `--use-package-json` to show only npm/yarn/pnpm scripts
- Use `--use-makefile` to show only Makefile targets

### UI Configuration Constants

These values are defined in source code and control the interface:

**Display sizing:**
```go
minViewportHeight     = 5     // Minimum terminal height
headerFooterLines     = 6     // Lines for title/filter/help
linesPerScriptOption  = 2     // Lines per script item
filterCharLimit       = 100   // Max filter input length
commandMaxWidthBuffer = 5     // Space for "..." truncation
```

**Keyboard controls:**
- **↑** - Move selection up (wraps around)
- **↓** - Move selection down (wraps around)
- **Enter** - Execute selected script
- **q** / **Ctrl+C** - Quit without running
- **Type** - Live filter scripts
- **Esc** - Clear filter
- **Backspace** - Delete filter character

### Database Schema

Location: `~/.config/alex-runner/alex-runner.sqlite.db`

**script_usage table:**
```sql
CREATE TABLE script_usage (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  directory TEXT NOT NULL,
  script_name TEXT NOT NULL,
  last_used TIMESTAMP NOT NULL,
  use_count INTEGER DEFAULT 1,
  UNIQUE(directory, script_name)
);

CREATE INDEX idx_directory ON script_usage(directory);
CREATE INDEX idx_frecency ON script_usage(directory, last_used DESC, use_count DESC);
```

**package_manager_cache table:**
```sql
CREATE TABLE package_manager_cache (
  directory TEXT PRIMARY KEY,
  package_manager TEXT NOT NULL,
  detected_at TIMESTAMP NOT NULL
);
```

### Time Display Format

Last used timestamps are displayed as:
- < 1 minute: "just now"
- < 1 hour: "N mins ago" or "1 min ago"
- < 24 hours: "Nh ago"
- < 7 days: "Nd ago"
- < 30 days: "Nw ago"
- ≥ 30 days: "Nmo ago"

## Tips

1. **Create an alias**: Add to your `.bashrc` or `.zshrc`:
   ```bash
   alias rr="alex-runner"
   alias rrl="alex-runner -l"
   ```

2. **Quick search with positional args**: Fastest way to run a script:
   ```bash
   rr build     # Search and run build script
   rr test      # Search and run test script
   rr lint      # Search and run lint script
   ```

3. **Pass arguments with aliases**: Works seamlessly with `--`:
   ```bash
   rr test -- --watch --coverage
   rrl test -- --reporter=verbose
   ```

4. **Reset history if habits change**: If your workflow changes, reset the history:
   ```bash
   alex-runner --reset
   ```

## Requirements

- Go 1.23 or higher
- A project with `package.json` and scripts defined, and/or a `Makefile` with targets

## Releases

alex-runner uses [Changesets](https://github.com/changesets/changesets) for version management and GitHub Actions + GoReleaser for building releases.

Pre-built binaries are available for:
- **macOS**: Intel (x86_64) and Apple Silicon (arm64)
- **Linux**: x86_64 and arm64
- **Windows**: x86_64

### Creating a Release

```bash
# 1. Add a changeset (during development)
pnpm changeset

# 2. Release (bumps version, updates changelog, commits, tags, and pushes)
pnpm release
```

GitHub Actions will automatically build and publish the binaries when the tag is pushed.

See [RELEASE.md](RELEASE.md) for more details.

## Contributing

Contributions are welcome! Please open an issue or PR.

## License

MIT

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Huh](https://github.com/charmbracelet/huh) - Interactive forms
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - Pure Go SQLite

Inspired by frecency algorithms in tools like Firefox and VS Code.
