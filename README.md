# Git Changelog Generator

A Go-based tool that generates formatted changelogs from Git tags and conventional commits.

## Features

- 📋 **Tag-based Changelog**: Lists all git tags sorted by time (newest to oldest)
- 📝 **Commit Grouping**: Groups commits by type using conventional commit format
- ⚠️ **Breaking Changes Detection**: Automatically detects and highlights breaking changes
- 🚫 **Commit Filtering**: Filters commits with ignore patterns (e.g., `[skip ci]`)
- 🎯 **Clean Output**: Shows only commit subject lines (removes descriptions)
- 📜 **Complete History**: Includes the oldest tag with all its historical commits

## Installation

```bash
go build
```

## Usage

```bash
./changelog
```

The tool will read the configuration from `.changelog.yml` and generate a changelog based on your Git tags.

## Configuration

Create a `.changelog.yml` file in your project root:

```yaml
git_path: "."  # Path to your git repository
ignore: "[skip ci]"  # Pattern to ignore commits
commit_groups:
  title_maps:
    breaking: Breaking Changes
    feat: Added
    fix: Fixed
    perf: Performance
    refactor: Changed
    docs: Documentation
    chore: Maintenance
    ci: CI/CD
    build: Build
    test: Testing
    style: Style
```

## Conventional Commits Format

The tool recognizes the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Examples

**Regular commit:**

```
feat: add user authentication
```

**Commit with scope:**

```
feat(auth): add OAuth2 support
```

**Breaking change with `!` indicator:**

```
feat!: remove deprecated API endpoints
```

**Breaking change with footer:**

```
feat: redesign user profile

BREAKING CHANGE: The user profile API has been completely redesigned.
Old endpoints are no longer supported.
```

## Breaking Changes Detection

The tool detects breaking changes in two ways:

### 1. Breaking Change Indicator (`!`)

Add an exclamation mark after the type/scope:

```
feat!: change API response format
fix(api)!: update authentication flow
```

### 2. Breaking Change Footer

Add a `BREAKING CHANGE:` or `BREAKING-CHANGE:` footer in the commit body:

```
feat: update database schema

BREAKING CHANGE: Database migration required. Run `migrate up` before deploying.
```

The breaking changes section will appear at the top of each tag's changelog with a ⚠️ icon.

## Output Format

```
========================================
v1.2.0
========================================

### ⚠️  Breaking Changes

  • abc1234 - **api**: change authentication flow to use JWT tokens

### Added

  • def5678 - **auth**: new OAuth2 provider support
  • ghi9012 - user profile page

### Fixed

  • jkl3456 - memory leak in cache
  • mno7890 - **db**: connection pool timeout

========================================
v1.1.0
========================================

### Added

  • pqr1234 - initial setup

========================================
v1.0.0 (oldest)
========================================

### Added

  • stu5678 - project initialization
```

## Project Structure

```
.
├── main.go      # Main logic for tag processing and changelog generation
├── commit.go    # Commit parsing and grouping logic
├── config.go    # Configuration management
├── tag.go       # Tag information struct
└── .changelog.yml  # Configuration file
```

## Tech Stack

- **[go-git/go-git](https://github.com/go-git/go-git)** - Git operations in pure Go
- **[gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)** - YAML parsing
- **Conventional Commits** - Commit message format

## How It Works

1. **Load Configuration**: Reads `.changelog.yml` for settings
2. **Fetch Tags**: Gets all tags from the repository and sorts by commit time
3. **Extract Commits**: For each tag pair, extracts commits between them
4. **Parse Commits**: Parses each commit using conventional commit format
5. **Detect Breaking Changes**: Checks for `!` indicator and `BREAKING CHANGE:` footer
6. **Group by Type**: Groups commits by their type (feat, fix, docs, etc.)
7. **Format Output**: Displays grouped commits with breaking changes at the top

## License

MIT
