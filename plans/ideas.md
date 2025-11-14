# alex-runner Feature Ideas & Improvements

This document contains potential features and improvements for alex-runner, organized by category and priority.

---

## 🎯 High-Value Features

### 1. ✅ Script Descriptions/Preview (ALREADY IMPLEMENTED)
Show the actual command being run before execution in the UI.

### 2. Success/Failure Tracking
**Status:** Proposed
**Priority:** High
**Effort:** Medium

Track exit codes and show reliability indicators:
```
> dev [★★★★★ 2h ago] ✓ 24/24 successful
> test [★★☆☆☆ 1d ago] ⚠ 3/10 failed
```

**Benefits:**
- Quick visual indicator for problematic scripts
- Could add `--show-failures` flag to see error history
- Helps identify flaky scripts

**Implementation Notes:**
- Add `exit_code` and `success_count`/`failure_count` columns to database
- Store exit code with each script execution
- Calculate success rate when displaying scripts

---

### 3. Execution Time Tracking
**Status:** Proposed
**Priority:** High
**Effort:** Medium

Show average/last runtime:
```
> build [★★★☆☆ 2h ago] ⏱ avg: 45s
> dev [★★★★★ 2h ago] ⏱ still running
```

**Benefits:**
- Helps set expectations for long-running scripts
- Could warn if script takes unusually long (performance regression detection)
- Useful for optimizing slow scripts

**Implementation Notes:**
- Add `duration_seconds` column to database
- Track start/end time of each execution
- Calculate average duration over last N runs
- Add optional `--show-duration` flag or always show if >5 seconds

---

### 4. Shell Completion
**Status:** Proposed
**Priority:** HIGH - Next to implement
**Effort:** Medium

Generate completion scripts for bash/zsh/fish:
```bash
# zsh example
alex-runner <TAB>  # Shows your most frecent scripts
alex-runner -l tes<TAB>  # Completes to 'test'
```

**Benefits:**
- Makes tool feel native and professional
- Faster workflow with tab completion
- Discovers flags and scripts without --list

**Implementation Notes:**
- Add `--generate-completion [bash|zsh|fish]` command
- Generate completion files for each shell
- Completions should be frecency-aware (suggest most frecent first)
- Include flag completions (-l, -s, --list, etc.)
- Installation instructions for each shell in README

**Technical Details:**
- Use Go's cobra/pflag if migrating from flags
- Or hand-craft completion scripts that call `alex-runner --list --json`
- Dynamic completions should query database for current directory
- Static completions for flags and subcommands

---

### 5. Workspace/Monorepo Support
**Status:** Proposed
**Priority:** High
**Effort:** High

```bash
alex-runner --workspace packages/frontend
alex-runner --all-workspaces  # Show scripts from all workspaces
```

**Benefits:**
- Essential for pnpm/yarn workspaces
- Could show workspace name in script list
- Run scripts across multiple packages

**Implementation Notes:**
- Parse workspace configuration from package.json
- Detect workspace root
- Allow filtering by workspace
- Show workspace name in script display

---

## 🔧 Quality of Life Improvements

### 6. Script Pinning/Favorites
**Status:** Proposed
**Priority:** Medium
**Effort:** Low

```bash
alex-runner --pin dev
# Now "dev" always appears first regardless of frecency
```

**Benefits:**
- Override frecency for critical scripts
- Useful for onboarding (pin common scripts for new developers)

**Implementation Notes:**
- Add `is_pinned` boolean to database
- Add `--pin <script>` and `--unpin <script>` flags
- Pinned scripts always appear first, then sorted by frecency
- Could have visual indicator (📌) in UI

---

### 7. Script Aliases
**Status:** Proposed
**Priority:** Medium
**Effort:** Medium

```bash
# In config file or flag
alex-runner --alias d=dev --alias b=build
alex-runner d  # Runs dev
```

**Benefits:**
- Faster typing for frequently used scripts
- Personal shortcuts

**Implementation Notes:**
- Store in config file (.alex-runner.json or similar)
- Resolve aliases before script lookup
- Show alias in UI alongside original name
- Allow global and per-project aliases

---

### 8. Recent History View
**Status:** Proposed
**Priority:** Medium
**Effort:** Low

```bash
alex-runner --history
# Shows last 10 runs with timestamps, duration, exit code
```

**Benefits:**
- Debug what was run recently
- Review execution history
- Useful for reporting/logging

**Implementation Notes:**
- Add `execution_history` table with timestamps
- Show last N executions with details
- Optional: add `--history <script-name>` to filter

---

### 9. Environment Variable Support
**Status:** Proposed
**Priority:** Medium
**Effort:** Medium

```bash
alex-runner dev --env NODE_ENV=production,DEBUG=true
# Or read from .env.local
```

**Benefits:**
- Quick environment switching
- No need to create separate scripts for different envs
- Integrate with existing .env files

**Implementation Notes:**
- Add `--env KEY=VALUE` flag (repeatable)
- Support `--env-file .env.local` to load file
- Merge with existing environment variables
- Could store commonly used env combos in config

---

### 10. Script Chaining
**Status:** Proposed
**Priority:** Medium
**Effort:** High

```bash
alex-runner -l lint,test,build  # Run in sequence
alex-runner -l lint+test+build  # Run in parallel (optional)
```

**Benefits:**
- Run common workflows quickly
- No need to create meta-scripts
- Short-circuit on failure

**Implementation Notes:**
- Parse comma-separated script names
- Run sequentially by default
- Stop on first failure (exit code != 0)
- Optional: support `+` for parallel execution
- Update frecency for all scripts in chain

---

## 🚀 Advanced Features

### 11. Watch Mode Integration
**Status:** Proposed
**Priority:** Low
**Effort:** High

```bash
alex-runner -w test  # Re-run test on file changes
# Uses native --watch if available, or implement own watcher
```

**Benefits:**
- Quick iteration during development
- No need to remember if script supports --watch

**Implementation Notes:**
- Check if script has native --watch support
- Otherwise implement file watcher
- Smart ignore patterns (node_modules, .git, etc.)
- Could use fsnotify or similar library

---

### 12. Output Filtering/Formatting
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

```bash
alex-runner --quiet  # Suppress verbose output
alex-runner --json   # JSON output for scripting
```

**Benefits:**
- Cleaner output for CI/CD
- Machine-readable format for scripts
- Focus on errors only

**Implementation Notes:**
- Add output level flags: --quiet, --verbose, --json
- For --json, output structured data about execution
- Preserve script's stdout/stderr handling

---

### 13. Config File Support
**Status:** Proposed
**Priority:** High (enables many other features)
**Effort:** Medium

```json
// .alex-runner.json or .alex-runner.toml
{
  "aliases": {"d": "dev", "b": "build"},
  "pinned": ["dev", "test"],
  "env": {"NODE_ENV": "development"},
  "frecency": {
    "frequencyWeight": 0.4,
    "recencyWeight": 0.6
  },
  "confirmBeforeRun": ["deploy*", "*:prod"],
  "theme": "dracula"
}
```

**Benefits:**
- Per-project configuration
- Team-wide settings via git
- Global user preferences
- Foundation for many other features

**Implementation Notes:**
- Support both JSON and TOML formats
- Search in: current dir, git root, home dir (~/.config/alex-runner/)
- Merge configs: project overrides global
- Use viper or similar config library
- Document all config options

---

### 14. Cross-Project Script Discovery
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

```bash
alex-runner --global-list
# Shows most-used scripts across ALL projects
# "Oh, I forgot we had a 'deploy:staging' script in this repo too!"
```

**Benefits:**
- Discover common patterns across projects
- Remember script names used elsewhere
- Insights into script naming conventions

**Implementation Notes:**
- Query database across all directories
- Show aggregated frecency scores
- Indicate which projects have each script
- Optional: suggest scripts from other projects

---

### 15. Script Suggestions Based on Context
**Status:** Proposed
**Priority:** Low
**Effort:** High

Smart suggestions based on git status:
- If you have uncommitted changes → suggest `test`, `lint`
- If on a feature branch → suggest `build`, `typecheck`
- If main branch → suggest `deploy`, `publish`

**Benefits:**
- Proactive workflow assistance
- Reduce errors (forgot to test before commit)
- Learn good practices

**Implementation Notes:**
- Integrate with git status checking
- Define rules in config file
- Show as hints/recommendations in UI
- Don't be too intrusive

---

### 16. Parallel Execution Mode
**Status:** Proposed
**Priority:** Low
**Effort:** High

```bash
alex-runner --parallel "lint,typecheck,test:unit"
# Runs all concurrently, shows live output for each
```

**Benefits:**
- Faster execution for independent scripts
- Useful for CI/CD optimization
- Better resource utilization

**Implementation Notes:**
- Run scripts in goroutines
- Multiplex output (with labels/colors)
- Collect all exit codes
- Fail if any script fails
- Consider adding to script chaining (#10)

---

## 📊 Analytics & Insights

### 17. Usage Statistics
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

```bash
alex-runner --stats
# Most used scripts this week/month
# Average runtime trends
# Success rate over time
```

**Benefits:**
- Understand workflow patterns
- Identify scripts that could be optimized
- Team insights (if shared database)

**Implementation Notes:**
- Add time-based queries to database
- Show top N scripts by period
- Graphs/charts in terminal (termui library)
- Optional: export to CSV/JSON

---

### 18. Script Dependency Graph
**Status:** Proposed
**Priority:** Low
**Effort:** High

Parse scripts that call other scripts:
```
build → typecheck → clean
deploy → build → test
```

**Benefits:**
- Understand script relationships
- Detect circular dependencies
- Optimize execution order

**Implementation Notes:**
- Parse script commands for npm/yarn/pnpm run calls
- Build dependency tree
- Visualize as ASCII tree or graph
- Could inform parallel execution strategy

---

### 19. Duplicate Script Detection
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

Warn if multiple scripts do similar things across projects or within same project.

**Benefits:**
- Reduce script bloat
- Identify candidates for consolidation
- Maintain consistency

**Implementation Notes:**
- Compare script commands (fuzzy matching)
- Warn about near-duplicates
- Suggest consolidation opportunities

---

## 🎨 UI/UX Enhancements

### 20. Themes/Customization
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

```bash
alex-runner --theme dracula
# Or in config file
```

**Benefits:**
- Personalization
- Better visibility in different terminal themes
- Accessibility (color blindness considerations)

**Implementation Notes:**
- Define color schemes using lipgloss
- Store in config file
- Preset themes: default, dracula, solarized, monokai, etc.
- Allow custom RGB values

---

### 21. Grouped Scripts Display
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

```
📦 Build & Deploy
  - build [★★★☆☆]
  - deploy [★☆☆☆☆]

🧪 Testing
  - test [★★★★☆]
  - test:e2e [★★☆☆☆]
```

**Benefits:**
- Better organization for large script lists
- Easier to scan visually
- Learn script categories

**Implementation Notes:**
- Group by prefix (test:*, build:*, deploy:*)
- Configurable groups in config file
- Collapsible groups in UI
- Show ungrouped scripts at end

---

### 22. Script Search Syntax Improvements
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

```bash
alex-runner "test NOT e2e"  # Boolean operators
alex-runner "build OR deploy"
alex-runner "test:*"  # Glob patterns
```

**Benefits:**
- More powerful search
- Better filtering for large script lists
- Familiar syntax for power users

**Implementation Notes:**
- Parse boolean operators (AND, OR, NOT)
- Support glob patterns with * and ?
- Combine with existing fuzzy search
- Document in help text

---

## 🔐 Safety Features

### 23. Dangerous Script Warnings
**Status:** Proposed
**Priority:** Medium
**Effort:** Low

```
⚠️  WARNING: This script contains 'rm -rf'
   > clean [★☆☆☆☆]
   Command: rm -rf dist node_modules
   Continue? [y/N]
```

**Benefits:**
- Prevent accidental data loss
- Extra caution for destructive operations
- Educational for junior developers

**Implementation Notes:**
- Scan script commands for dangerous patterns
- Patterns: rm -rf, dd, mkfs, DROP TABLE, etc.
- Configurable warning patterns
- Require explicit confirmation
- Could disable with --force flag

---

### 24. Dry Run Mode
**Status:** Proposed
**Priority:** Medium
**Effort:** Low

```bash
alex-runner --dry-run build
# Shows what would be executed without running it
```

**Benefits:**
- Preview before execution
- Test script selection
- Safe exploration

**Implementation Notes:**
- Show full command that would be executed
- Show environment variables that would be set
- Don't update frecency database
- Simple flag check before execution

---

### 25. Script Permissions/Confirmations
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

Require confirmation for scripts matching patterns:
```json
{
  "confirmBeforeRun": ["deploy*", "*:prod", "publish"]
}
```

**Benefits:**
- Prevent accidental production deploys
- Extra safety for critical scripts
- Configurable per project/team

**Implementation Notes:**
- Check script name against patterns in config
- Show confirmation prompt before execution
- Related to #23 (dangerous script warnings)
- Could combine both features

---

## 🔄 Integration Features

### 26. Git Hook Integration
**Status:** Proposed
**Priority:** Low
**Effort:** Medium

```bash
alex-runner --install-hooks
# Suggests running lint/test before commits
```

**Benefits:**
- Automated quality checks
- Consistent pre-commit workflow
- Team standards enforcement

**Implementation Notes:**
- Generate git hooks in .git/hooks/
- Run most frecent test/lint scripts
- Fast-fail on errors
- Respect existing hooks (append, don't replace)
- Consider using husky or similar if present

---

### 27. CI/CD Export
**Status:** Proposed
**Priority:** Low
**Effort:** High

```bash
alex-runner --export-github-actions
# Generates workflow with your most-used scripts
```

**Benefits:**
- Quick CI/CD setup
- Based on actual usage patterns
- Reduces CI/CD configuration time

**Implementation Notes:**
- Query most frecent scripts
- Generate YAML for GitHub Actions, GitLab CI, etc.
- Include test, build, deploy scripts
- Template-based generation
- Validate generated config

---

## 🎯 Top Priority Recommendations

Based on value vs. effort and user feedback:

1. **Shell Completion** (#4) - HIGH PRIORITY, NEXT TO IMPLEMENT
   - Makes tool feel professional and native
   - Medium effort, high value
   - Benefits all users immediately

2. **Config File Support** (#13)
   - Enables many other features
   - Foundation for customization
   - Medium effort, enables high value

3. **Execution Time Tracking** (#3)
   - Very useful for all users
   - Medium effort, high value
   - Helps identify slow scripts

4. **Success/Failure Tracking** (#2)
   - Identifies problematic scripts
   - Medium effort, high value
   - Useful for debugging

5. **Workspace Support** (#5)
   - Critical for monorepo users
   - High effort, high value for target users
   - Growing importance with monorepo trend

---

## Ideas Considered But Rejected

### Task Runner Integration (Taskfile, Just, etc.)
**Reason:** User doesn't use these tools, out of scope for project focus on npm/yarn/pnpm scripts and Makefiles.

---

## Contributing Ideas

Have a feature idea? Consider:
- **Value:** How many users will benefit?
- **Effort:** How complex is the implementation?
- **Scope:** Does it fit the project's focus?
- **Alternatives:** Could this be solved differently?

Open an issue or PR to discuss!
