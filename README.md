# alex-runner

A smart, frecency-based npm script runner that learns which scripts you use most often and suggests them automatically.

## problem

Sometimes there's a lot of scripts but you only really use a few of them and forget their names. alex-runner helps you remember which scripts you use most often and suggests them automatically.

## Features

- **Frecency-based suggestions**: Combines frequency and recency to suggest the scripts you're most likely to need
- **Live filtering**: Type to instantly filter scripts - no special keys needed
- **Beautiful TUI**: Powered by Bubble Tea with syntax highlighting and clear command previews
- **Multi-package manager**: Automatically detects npm, pnpm, or yarn
- **Makefile support**: Run Makefile targets alongside npm scripts
- **Per-directory tracking**: Each project has its own usage history
- **Fuzzy search**: Quickly find scripts by name or command content
- **Smart search ranking**: 6-tier priority system from exact matches to fuzzy command matches
- **Zero configuration**: Just install and run

## Demo Flow

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
go install github.com/alexanderchan/alex-runner@latest
```

for development

```bash
go install ./cmd/alex-runner
```

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

### Add to PATH

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

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
| `--reset` | | boolean | false | Clear usage history for current directory |
| `--global-reset` | | boolean | false | Clear all usage history |
| `--use-package-json` | | boolean | false | Only show package.json scripts (ignore Makefile) |
| `--use-makefile` | | boolean | false | Only show Makefile targets (ignore package.json) |
| `--no-cache` | | boolean | false | Re-detect package manager instead of using cached detection |
| `--help` | `-h` | boolean | false | Show help message |
| (positional arg) | | string | "" | Same as `--search` - `alex-runner build` |

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
- **↑** / **k** - Move selection up (wraps around)
- **↓** / **j** - Move selection down (wraps around)
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

3. **Reset history if habits change**: If your workflow changes, reset the history:
   ```bash
   alex-runner --reset
   ```

## Requirements

- Go 1.23 or higher
- A project with `package.json` and scripts defined, and/or a `Makefile` with targets

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
