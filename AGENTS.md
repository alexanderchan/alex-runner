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

### Lucky mode
```bash
alex-runner -l              # Run most frecent script immediately
alex-runner -l build        # Run first "build" match immediately
```

## Implementation Notes

### Fuzzy Matching
- The fuzzy matching is powered by `github.com/lithammer/fuzzysearch/fuzzy`
- Located in [internal/search.go](internal/search.go)
- The `SearchScripts` function provides ranking:
  - Exact match: 1000
  - Prefix match: 500
  - Contains: 300
  - Fuzzy name match: 200
  - Command contains: 100
  - Fuzzy command match: 50

### UI Filtering
- The interactive UI uses the same `SearchScripts` function (as of the recent fix)
- Previously used simple substring matching, now consistent with fuzzy search
- Located in [internal/ui.go:353](internal/ui.go#L353)

## Testing
Run tests with:
```bash
go test ./internal -v
```

Specific fuzzy matching test:
```bash
go test ./internal -run TestSearchScriptsFuzzyMatch -v
```
