# Release Process

This document describes how to create a new release of alex-runner.

## Prerequisites

- Push access to the repository
- pnpm installed (for changesets)
- GoReleaser installed locally (optional, for testing)

## Release Workflow

alex-runner uses [Changesets](https://github.com/changesets/changesets) for version management and GitHub Actions + GoReleaser for building releases.

### 1. Add a Changeset (during development)

When you make a change that should be released:

```bash
pnpm changeset
```

This will prompt you to:
- Select the type of change (major/minor/patch)
- Write a summary of the change

A markdown file will be created in `.changeset/` - commit this with your PR.

### 2. Version and Release

When ready to release:

```bash
# Consume all changesets, bump version, update CHANGELOG.md
pnpm changeset version

# Review the changes
git diff

# Commit
git add .
git commit -m "chore: release v0.3.0"
git push origin main

# Create and push the tag
git tag v0.3.0
git push origin v0.3.0
```

### 3. CI Builds the Release

When the tag is pushed, GitHub Actions will:
- Build binaries for all supported platforms
- Create archives (`.tar.gz` for Unix, `.zip` for Windows)
- Generate checksums
- Create a GitHub release with all artifacts

## Testing Locally

You can test the release process locally before pushing a tag:

### Install GoReleaser

```bash
# macOS
brew install goreleaser

# Or using Go
go install github.com/goreleaser/goreleaser@latest
```

### Test the Build

```bash
# Build without releasing (snapshot mode)
goreleaser build --snapshot --clean

# Check the dist/ directory for built binaries
ls -la dist/
```

### Test the Full Release Process

```bash
# Dry run - builds everything but doesn't publish
goreleaser release --snapshot --clean

# Check the generated archives
ls -la dist/*.tar.gz dist/*.zip
```

## Supported Platforms

The release process builds binaries for:

- **macOS**:
  - `alex-runner_Darwin_x86_64.tar.gz` (Intel)
  - `alex-runner_Darwin_arm64.tar.gz` (Apple Silicon)

- **Linux**:
  - `alex-runner_Linux_x86_64.tar.gz` (AMD64)
  - `alex-runner_Linux_arm64.tar.gz` (ARM64)

- **Windows**:
  - `alex-runner_Windows_x86_64.zip` (AMD64)

## Troubleshooting

### Build Fails

If the build fails:
1. Check the GitHub Actions logs
2. Test locally with `goreleaser build --snapshot --clean`
3. Ensure all dependencies are properly declared in `go.mod`

### Release Not Created

If the release isn't created:
1. Verify the tag follows the `v*` pattern (e.g., `v0.3.0`)
2. Check that the GitHub Actions workflow has write permissions
3. Ensure `GITHUB_TOKEN` has the necessary permissions

## Notes

- alex-runner uses `modernc.org/sqlite`, which is a pure Go implementation
- No CGO required, making cross-compilation simple
- All builds are static binaries with no external dependencies

