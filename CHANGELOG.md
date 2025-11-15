# Changelog

All notable changes to the Git Changelog Generator project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### ⚠️ Breaking Changes

- Refactored commit-related functions into separate `commit.go` module

### Added

- **Breaking Changes Detection**: Automatically detects and highlights breaking changes in changelogs
  - Detects `!` indicator after commit type (e.g., `feat!:` or `feat(scope)!:`)
  - Detects `BREAKING CHANGE:` or `BREAKING-CHANGE:` footer in commit body
  - Displays breaking changes at the top of each tag section with ⚠️ icon
- **CommitInfo Struct**: New struct to hold parsed commit information including:
  - Commit type and scope
  - Subject and body
  - Breaking change status and description
- **Enhanced Commit Parsing**: `parseCommit()` function extracts all conventional commit components
- **Scope Support**: Commit scopes are now displayed in bold (e.g., `**api**: change format`)
- Configuration option for breaking changes title in `.changelog.yml`
- Comprehensive documentation:
  - `README.md` with full feature documentation
  - `EXAMPLES.md` with real-world usage examples
  - Migration guide for breaking changes

### Changed

- Moved all commit-related functions from `main.go` to new `commit.go` file:
  - `getFirstLine()`
  - `extractCommitType()`
  - `shouldIgnoreCommit()`
  - `printGroupedCommits()`
  - `parseCommit()` (new)
- Improved regex patterns to handle optional scopes and breaking change indicators
- Enhanced commit message parsing to extract body and footers
- Updated `.changelog.yml` to include `breaking: Breaking Changes` in title_maps

### Fixed

- Cleaned up duplicate code in main.go
- Improved error handling in `getCommitsBetween()`

## Project Structure

```
changelog/
├── main.go           # Main logic for tag processing
├── commit.go         # Commit parsing and grouping (NEW)
├── config.go         # Configuration management
├── tag.go            # Tag information struct
├── .changelog.yml    # Configuration file
├── README.md         # Project documentation
├── EXAMPLES.md       # Usage examples (NEW)
└── CHANGELOG.md      # This file (NEW)
```

## Technical Details

### Breaking Change Detection Algorithm

1. Parse commit message into subject and body
2. Check subject for `!` indicator using regex: `^(\w+)(?:\(([^)]+)\))?(!)?:\s*(.*)$`
3. Check body for `BREAKING CHANGE:` or `BREAKING-CHANGE:` footer
4. Extract breaking change description from footer if present
5. Mark commit as breaking if either indicator is found

### Output Format

```
========================================
v2.0.0
========================================

### ⚠️  Breaking Changes

  • abc1234 - **scope**: breaking change description

### Added

  • def5678 - **scope**: new feature
```

## Migration Guide

If you're upgrading from a previous version:

1. **Update Configuration**: Add `breaking: Breaking Changes` to your `.changelog.yml`:
   ```yaml
   commit_groups:
     title_maps:
       breaking: Breaking Changes  # Add this line
       feat: Added
       fix: Fixed
       # ... other types
   ```

2. **Rebuild**: Run `go build` to compile the new version

3. **No Breaking Changes**: This update is backward compatible. Existing changelogs will continue to work, with the addition of breaking changes detection.

## Contributing

When contributing to this project, please follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- Use `feat:` for new features
- Use `fix:` for bug fixes
- Add `!` after the type for breaking changes (e.g., `feat!:`)
- Add `BREAKING CHANGE:` footer for detailed breaking change descriptions
- Use scopes to indicate the affected module (e.g., `feat(parser):`)

## License

MIT