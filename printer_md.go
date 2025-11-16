package main

import (
	"fmt"

	"github.com/go-git/go-git/v6"
)

type Commit struct {
	Hash        string
	Title       string // feat, fix, refactor, etc...
	Scope       string
	Description string
	Author      string
	IsBreaking  bool
}

type Tag struct {
	Name    string
	Commits []Commit
}

type ChangelogData struct {
	Tags []Tag
}

type MarkdownPrinter struct {
	config *Config
	repo   *git.Repository

	Data ChangelogData
}

func (c *MarkdownPrinter) MapData(tags []*TagInfo) {
	// Process commits between consecutive tags
	for i := 0; i < len(tags)-1; i++ {
		newerTag := tags[i]
		olderTag := tags[i+1]

		// Get commits between the two tags
		commits, err := getCommitsBetween(c.repo, &olderTag.Commit.Hash, newerTag.Commit.Hash)
		if err != nil {
			panic(err)
		}

		tagData := Tag{
			Name:    newerTag.Name,
			Commits: []Commit{},
		}

		// Parse and filter commits
		for _, commit := range commits {
			if shouldIgnoreCommit(commit.Message, c.config.Ignore) {
				continue
			}

			info := parseCommit(commit)
			tagData.Commits = append(tagData.Commits, Commit{
				Hash:        info.Commit.Hash.String()[:7],
				Title:       info.Type,
				Scope:       info.Scope,
				Description: info.Subject,
				Author:      info.Commit.Author.Name,
				IsBreaking:  info.IsBreaking,
			})
		}

		c.Data.Tags = append(c.Data.Tags, tagData)
	}

	// Process commits for the oldest tag (all commits up to and including that tag)
	if len(tags) > 0 {
		oldestTag := tags[len(tags)-1]

		// Get commits between the two tags
		commits, err := getCommitsBetween(c.repo, nil, oldestTag.Commit.Hash)
		if err != nil {
			panic(err)
		}

		tagData := Tag{
			Name:    oldestTag.Name,
			Commits: []Commit{},
		}

		// Parse and filter commits
		for _, commit := range commits {
			if shouldIgnoreCommit(commit.Message, c.config.Ignore) {
				continue
			}

			info := parseCommit(commit)
			tagData.Commits = append(tagData.Commits, Commit{
				Hash:        info.Commit.Hash.String()[:7],
				Title:       info.Type,
				Scope:       info.Scope,
				Description: info.Subject,
				Author:      info.Commit.Author.Name,
				IsBreaking:  info.IsBreaking,
			})
		}

		c.Data.Tags = append(c.Data.Tags, tagData)
	}
}

// escapePipes escapes pipe characters in strings to prevent breaking markdown tables
func escapePipes(s string) string {
	return s
}

// Print implements Printer.
func (c *MarkdownPrinter) Print(current string) {
	// Print each tag section
	for _, tag := range c.Data.Tags {
		title := fmt.Sprintf("\n## %s", tag.Name)
		if current == tag.Name {
			title = fmt.Sprintf("%s - Current Release", title)
		}
		fmt.Printf("%s\n\n", title)

		if len(tag.Commits) == 0 {
			fmt.Println("_No commits between these tags_")
			fmt.Println()
			continue
		}

		// Group commits by type
		groups := make(map[string][]Commit)
		breakingChanges := make([]Commit, 0)

		for _, commit := range tag.Commits {
			// Check for breaking changes
			if commit.IsBreaking {
				breakingChanges = append(breakingChanges, commit)
			} else {
				groups[commit.Title] = append(groups[commit.Title], commit)
			}
		}

		// Define order of groups based on config
		var groupOrder []string
		for key := range c.config.CommitGroups.TitleMaps {
			if key == "breaking" {
				continue
			}
			if _, exists := groups[key]; exists {
				groupOrder = append(groupOrder, key)
			}
		}
		// Add "other" category at the end if it exists
		if _, exists := groups["other"]; exists {
			groupOrder = append(groupOrder, "other")
		}

		// Print breaking changes first if any
		if len(breakingChanges) > 0 {
			breakingTitle := c.config.CommitGroups.TitleMaps["breaking"]
			if breakingTitle == "" {
				breakingTitle = "Breaking Changes"
			}

			fmt.Printf("### ⚠️  %s\n\n", breakingTitle)
			fmt.Println("| Commit | Scope | Description | Author |")
			fmt.Println("|--------|-------|-------------|--------|")
			for _, commit := range breakingChanges {
				scope := "-"
				if commit.Scope != "" {
					scope = commit.Scope
				}
				fmt.Printf("| `%s` | %s | %s | %s |\n", commit.Hash, scope, escapePipes(commit.Description), commit.Author)
			}
			fmt.Println()
		}

		// Print each group
		for _, groupType := range groupOrder {
			title := c.config.CommitGroups.TitleMaps[groupType]
			if title == "" {
				title = "Other"
			}

			fmt.Printf("### %s\n\n", title)
			fmt.Println("| Commit | Scope | Description | Author |")
			fmt.Println("|--------|-------|-------------|--------|")
			for _, commit := range groups[groupType] {
				scope := "-"
				if commit.Scope != "" {
					scope = commit.Scope
				}
				fmt.Printf("| `%s` | %s | %s | %s |\n", commit.Hash, scope, escapePipes(commit.Description), commit.Author)
			}
			fmt.Println()
		}
	}
}

func NewMarkdownPrinter(config *Config, repo *git.Repository) Printer {
	return &MarkdownPrinter{
		config: config,
		repo:   repo,
	}
}
