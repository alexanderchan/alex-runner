# Alex Runner - Agent Notes

## Useful Commands

### List all scripts with frecency scores
```bash
alex-runner --list
```

This command shows all available scripts from both Makefile and package.json (if present), sorted by frecency scores. It displays:
- Script names
- Commands
- Source (make/npm/pnpm/yarn)
- Usage statistics

Useful for debugging and understanding which scripts are available and which ones are used most frequently.

### Search and filter
```bash
alex-runner -s <search-term>
# or
alex-runner <search-term>
```

The interactive filter supports **fuzzy matching**, so queries like "dmo" will match "demo-record", "demo-generate", etc.

#### Keyboard Shortcuts in Interactive Mode
- `↑`/`↓` or `j`/`k` - Navigate through results
- `Enter` - Select and run script
- `q` or `Ctrl+C` - Quit without running

**Clear filter (show all scripts) - multiple options:**
- `Esc` - Universal, works everywhere
- `Ctrl+U` - Standard terminal "clear line" shortcut
- `Alt+Backspace` - May work as `Cmd+Backspace` on Mac terminals
- `Ctrl+Backspace` - Works on some terminal emulators
- `Ctrl+W` - Standard "delete word backward"

All five clear shortcuts are supported to accommodate different terminal preferences and platform behaviors.

### Lucky mode
```bash
alex-runner -l              # Run most frecent script immediately
alex-runner -l build        # Run first "build" match immediately
```

## Implementation Notes

### Fuzzy Matching
The search system uses a **hybrid approach** for optimal results:

#### Single-Word Queries (e.g., "dmo", "build", "test")
- Powered by `github.com/lithammer/fuzzysearch/fuzzy`
- Located in [internal/search.go](internal/search.go)
- Ranking system:
  - Exact match: 1000
  - Prefix match: 500
  - Contains: 300
  - Fuzzy name match: 200
  - Command contains: 100
  - Fuzzy command match: 50

Example: `dmo` matches `demo-record`, `demo-generate`

#### Multi-Word Queries (e.g., "docker build", "build docker")
- Powered by `github.com/schollz/closestmatch` (bag-of-words approach)
- **Order-independent**: "docker build" and "build docker" both match "start-docker:traefik:build"
- Matches words anywhere in the script name or command
- Uses n-gram matching (sizes 2, 3, 4) for fuzzy word matching

Example:
- `docker build` → matches `start-docker:traefik:build` ✓
- `build docker` → matches `start-docker:traefik:build` ✓
- `traefik build` → matches `start-docker:traefik:build` ✓

#### Implementation Details
- Single vs multi-word detection happens in `SearchScripts()` using `strings.Fields()`
- Multi-word queries are handled by `searchWithClosestMatch()` function
- The hybrid approach gives you the best of both:
  - Fast fuzzy matching for abbreviations
  - Flexible word-order-independent matching for complex queries

### UI Filtering
- The interactive UI uses the same `SearchScripts` function for consistency
- Located in [internal/ui.go:362](internal/ui.go#L362)
- Both CLI search (`-s` flag) and interactive filtering use identical matching logic

## Testing
Run tests with:
```bash
go test ./internal -v
```

Specific fuzzy matching test:
```bash
go test ./internal -run TestSearchScriptsFuzzyMatch -v
```
