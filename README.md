# alex-runner

A smart, frecency-based npm script runner that learns which scripts you use most often and suggests them automatically.

## problem

Sometimes there's a lot of scripts but you only really use a few of them and forget their names. alex-runner helps you remember which scripts you use most often and suggests them automatically.

## Features

- **Frecency-based suggestions**: Combines frequency and recency to suggest the scripts you're most likely to need
- **Live filtering**: Type to instantly filter scripts - no special keys needed
- **Beautiful TUI**: Powered by Bubble Tea with syntax highlighting and clear command previews
- **Multi-package manager**: Automatically detects npm, pnpm, or yarn
- **Per-directory tracking**: Each project has its own usage history
- **Fuzzy search**: Quickly find scripts by name or command content
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

alex-runner automatically detects your package manager:

- Looks for `pnpm-lock.yaml` → uses `pnpm`
- Looks for `yarn.lock` → uses `yarn`
- Looks for `package-lock.json` → uses `npm`
- Defaults to `npm`

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

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--last` | `-l` | "I'm feeling lucky" - run most frecent immediately |
| `--search` | `-s` | Show selector filtered to search term |
| `--list` | | List all scripts with frecency scores |
| `--reset` | | Clear usage history for current directory |
| `--global-reset` | | Clear all usage history |
| `--help` | `-h` | Show help message |
| (positional arg) | | Same as `--search` - `alex-runner build` |

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
- A project with `package.json` and scripts defined

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
