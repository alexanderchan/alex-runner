# Changelog

## 0.2.4

### Patch Changes

- ff7e428: just testing again

## 0.2.3

### Patch Changes

- 93eab55: test release script

## 0.2.2

### Patch Changes

- testing release

## 0.2.1

### Patch Changes

- 5b96d0e: initial patch release test

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-01-XX

### Added

- **Script Pinning**: Pin your most important scripts to always appear first
  - Use `--pin <script>` to pin a script from the command line
  - Use `--unpin <script>` to unpin a script
  - Press `alt-p` (or `option-p` on Mac) in the interactive UI to toggle pin status
  - Pinned scripts show a 📌 indicator
  - When multiple sources (Makefile and package.json) have the same script name, you can:
    - Choose which one to pin (select `1`, `2`, etc.)
    - Pin all of them at once (select `all`)
- Source-aware script tracking: Scripts from Makefile and package.json are now tracked separately
  - Fixes issue where identically named scripts from different sources would conflict
  - Each script+source combination has its own frecency score and pin status

### Changed

- Database schema updated to include `source` field for better script disambiguation
- Scripts are now uniquely identified by `(directory, script_name, source)` composite key
- Pinned scripts always appear first in the list, sorted by frecency among themselves

### Breaking Changes

- **Database schema change**: The database structure has been updated to track script sources
- **Migration required**: If you're upgrading from v0.1.0, you'll need to clear your database:
  ```bash
  rm ~/.config/alex-runner/alex-runner.sqlite.db
  ```
- After removing the database, alex-runner will automatically create a new one with the updated schema
- Your usage history will be reset, but the tool will start tracking with the new schema

### Migration Notes

The database now uses a composite key `(directory, script_name, source)` instead of just `(directory, script_name)`. This means:

- Scripts from Makefile (source: "make") and package.json (source: "npm"/"pnpm"/"yarn") are tracked separately
- Pin status is per script+source combination
- Frecency scores are tracked independently for each source

If you have critical usage history you want to preserve, please stay on v0.1.0. Otherwise, the fresh start with improved tracking is recommended.

## [0.1.0] - 2025-01-XX

### Added

- Initial release
- Frecency-based script selection for npm/pnpm/yarn/Makefile scripts
- Interactive script selector with live filtering
- Usage tracking per directory
- Support for both Makefile targets and package.json scripts
- Shell completion generation (bash, zsh, fish)
- "I'm feeling lucky" mode with `-l` flag
- Search functionality with `-s` flag
- Package manager auto-detection with caching
