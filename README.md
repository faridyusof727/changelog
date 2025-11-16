# Git Changelog Generator

A Go-based tool that generates formatted changelogs from Git tags and conventional commits, outputting in Markdown table format.

## Features

- 📋 **Tag-based Changelog**: Lists all git tags sorted by time (newest to oldest)
- 📝 **Commit Grouping**: Groups commits by type using conventional commit format
- ⚠️ **Breaking Changes Detection**: Automatically detects and highlights breaking changes
- 🚫 **Commit Filtering**: Filters commits with ignore patterns (e.g., `[skip ci]`)
- 📊 **Markdown Tables**: Outputs commits in clean, readable markdown tables
- 👤 **Author Tracking**: Displays commit author for each change
- 🔍 **Scope Support**: Shows commit scope when available
- 📜 **Complete History**: Includes the oldest tag with all its historical commits

## Installation

### From Source

```bash
git clone https://github.com/faridyusof727/changelog.git
cd changelog
go build
```

### Using Go Install

```bash
go install github.com/faridyusof727/changelog@latest
```

## Usage

```bash
./changelog > CHANGELOG.md
```

The tool will read the configuration from `.changelog.yml` and generate a changelog based on your Git tags, outputting to stdout.

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

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Examples

**Regular commit:**

```text
feat: add user authentication
```

**Commit with scope:**

```text
feat(auth): add OAuth2 support
```

**Breaking change with `!` indicator:**

```text
feat!: remove deprecated API endpoints
```

**Breaking change with footer:**

```text
feat: redesign user profile

BREAKING CHANGE: The user profile API has been completely redesigned.
Old endpoints are no longer supported.
```

## Breaking Changes Detection

The tool detects breaking changes in two ways:

### 1. Breaking Change Indicator (`!`)

Add an exclamation mark after the type/scope:

```text
feat!: change API response format
fix(api)!: update authentication flow
```

### 2. Breaking Change Footer

Add a `BREAKING CHANGE:` or `BREAKING-CHANGE:` footer in the commit body:

```text
feat: update database schema

BREAKING CHANGE: Database migration required. Run `migrate up` before deploying.
```

The breaking changes section will appear at the top of each tag's changelog with a ⚠️ icon.

## Output Format

The tool generates markdown-formatted changelogs with tables for each commit group:

```markdown
## v1.2.0

### ⚠️  Breaking Changes

| Commit | Scope | Description | Author |
|--------|-------|-------------|--------|
| `abc1234` | api | change authentication flow to use JWT tokens | John Doe |

### Added

| Commit | Scope | Description | Author |
|--------|-------|-------------|--------|
| `def5678` | auth | new OAuth2 provider support | Jane Smith |
| `ghi9012` | - | user profile page | John Doe |

### Fixed

| Commit | Scope | Description | Author |
|--------|-------|-------------|--------|
| `jkl3456` | - | memory leak in cache | Bob Johnson |
| `mno7890` | db | connection pool timeout | Jane Smith |

## v1.1.0

### Added

| Commit | Scope | Description | Author |
|--------|-------|-------------|--------|
| `pqr1234` | - | initial setup | John Doe |

## v1.0.0 (oldest)

### Added

| Commit | Scope | Description | Author |
|--------|-------|-------------|--------|
| `stu5678` | - | project initialization | John Doe |
```

## Project Structure

```text
.
├── main.go           # Entry point and main orchestration
├── config.go         # Configuration file parsing (YAML)
├── tag.go            # Tag loading and sorting logic
├── commit.go         # Commit parsing and conventional commit handling
├── printer.go        # Printer interface definition
├── printer_md.go     # Markdown table printer implementation
├── .changelog.yml    # Configuration file (user-created)
└── EXAMPLES.md       # Detailed examples of commit formats
```

## Tech Stack

- **[go-git/go-git](https://github.com/go-git/go-git)** - Git operations in pure Go
- **[gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)** - YAML parsing
- **Conventional Commits** - Commit message format

## How It Works

1. **Load Configuration**: Reads `.changelog.yml` for settings and commit group mappings
2. **Open Repository**: Opens the Git repository at the specified path
3. **Fetch Tags**: Gets all tags from the repository and sorts by commit time (newest first)
4. **Load Tag Info**: Resolves each tag to its commit (handles both lightweight and annotated tags)
5. **Extract Commits**: For each consecutive tag pair, extracts commits between them
6. **Parse Commits**: Parses each commit using conventional commit format:
   - Extracts type, scope, and subject from commit message
   - Detects breaking changes via `!` indicator or `BREAKING CHANGE:` footer
   - Filters out commits matching the ignore pattern
7. **Group by Type**: Groups commits by their type (feat, fix, docs, etc.)
8. **Format Output**: Generates markdown tables with:
   - Breaking changes section first (if any)
   - Commit groups in configured order
   - Each row showing: commit hash, scope, description, and author
  
## TODO

- [ ] Add HTML static file printer for web-based changelog viewing
- [ ] Add JSON printer for programmatic consumption
- [ ] Add plain text printer for simple output
- [ ] Support custom templates for different output formats
- [ ] Add CLI flags for output format selection

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT
