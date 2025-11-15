package main

import (
	"fmt"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func Print(config *Config, repo *git.Repository, tags []*TagInfo) {
	// Print commits between consecutive tags
	for i := 0; i < len(tags)-1; i++ {
		newerTag := tags[i]
		olderTag := tags[i+1]

		fmt.Printf("\n========================================\n")
		fmt.Printf("%s\n", newerTag.Name)
		fmt.Printf("========================================\n\n")

		// Get commits between the two tags
		commits, err := getCommitsBetween(repo, &olderTag.Commit.Hash, newerTag.Commit.Hash)
		if err != nil {
			panic(err)
		}

		if len(commits) == 0 {
			fmt.Println("  No commits between these tags")
		} else {
			printGroupedCommits(commits, config)
		}
	}

	// Print commits for the oldest tag (all commits up to and including that tag)
	if len(tags) > 0 {
		oldestTag := tags[len(tags)-1]

		fmt.Printf("\n========================================\n")
		fmt.Printf("%s (oldest)\n", oldestTag.Name)
		fmt.Printf("========================================\n\n")

		// Get commits between the two tags
		commits, err := getCommitsBetween(repo, nil, oldestTag.Commit.Hash)
		if err != nil {
			panic(err)
		}

		if len(commits) == 0 {
			fmt.Println("  No commits between these tags")
		} else {
			printGroupedCommits(commits, config)
		}
	}
}

// PrintMD prints changelog information in markdown format with tables for grouped commits
func PrintMD(config *Config, repo *git.Repository, tags []*TagInfo) {
	// Print commits between consecutive tags
	for i := 0; i < len(tags)-1; i++ {
		newerTag := tags[i]
		olderTag := tags[i+1]

		fmt.Printf("\n## %s\n\n", newerTag.Name)

		// Get commits between the two tags
		commits, err := getCommitsBetween(repo, &olderTag.Commit.Hash, newerTag.Commit.Hash)
		if err != nil {
			panic(err)
		}

		if len(commits) == 0 {
			fmt.Println("_No commits between these tags_")
			fmt.Println()
		} else {
			printGroupedCommitsInTableForMarkdown(commits, config)
		}
	}

	// Print commits for the oldest tag (all commits up to and including that tag)
	if len(tags) > 0 {
		oldestTag := tags[len(tags)-1]

		fmt.Printf("\n## %s (oldest)\n\n", oldestTag.Name)

		// Get commits between the two tags
		commits, err := getCommitsBetween(repo, nil, oldestTag.Commit.Hash)
		if err != nil {
			panic(err)
		}

		if len(commits) == 0 {
			fmt.Println("_No commits between these tags_")
			fmt.Println()
		} else {
			printGroupedCommitsInTableForMarkdown(commits, config)
		}
	}
}

// printGroupedCommits groups commits by type and prints them
func printGroupedCommits(commits []*object.Commit, config *Config) {
	// Parse and filter commits
	var commitInfos []*CommitInfo
	var breakingChanges []*CommitInfo

	for _, commit := range commits {
		if shouldIgnoreCommit(commit.Message, config.Ignore) {
			continue
		}

		info := parseCommit(commit)
		commitInfos = append(commitInfos, info)

		if info.IsBreaking {
			breakingChanges = append(breakingChanges, info)
		}
	}

	// If all commits were filtered out, show message
	if len(commitInfos) == 0 {
		fmt.Println("  No commits (all filtered)")
		return
	}

	// Print breaking changes first if any exist
	if len(breakingChanges) > 0 {
		breakingTitle := config.CommitGroups.TitleMaps["breaking"]
		if breakingTitle == "" {
			breakingTitle = "Breaking Changes"
		}

		fmt.Printf("### ⚠️  %s\n\n", breakingTitle)
		for _, info := range breakingChanges {
			message := info.Subject
			if info.BreakingMessage != "" {
				message = info.BreakingMessage
			}

			scope := ""
			if info.Scope != "" {
				scope = fmt.Sprintf("**%s**: ", info.Scope)
			}

			fmt.Printf("  • %s - %s%s\n",
				info.Commit.Hash.String()[:7],
				scope,
				message)
		}
		fmt.Println()
	}

	// Group non-breaking commits by type
	groups := make(map[string][]*CommitInfo)
	for _, info := range commitInfos {
		groups[info.Type] = append(groups[info.Type], info)
	}

	// Define order of groups based on config
	var groupOrder []string
	for key := range config.CommitGroups.TitleMaps {
		if key == "breaking" {
			continue // Already printed above
		}
		if _, exists := groups[key]; exists {
			groupOrder = append(groupOrder, key)
		}
	}
	// Add "other" category at the end if it exists
	if _, exists := groups["other"]; exists {
		groupOrder = append(groupOrder, "other")
	}

	// Print each group
	for _, groupType := range groupOrder {
		title := config.CommitGroups.TitleMaps[groupType]
		if title == "" {
			title = "Other"
		}

		fmt.Printf("### %s\n\n", title)
		for _, info := range groups[groupType] {
			message := info.Subject

			scope := ""
			if info.Scope != "" {
				scope = fmt.Sprintf("**%s**: ", info.Scope)
			}

			fmt.Printf("  • %s - %s%s\n",
				info.Commit.Hash.String()[:7],
				scope,
				message)
		}
		fmt.Println()
	}
}

// printGroupedCommitsInTableForMarkdown groups commits by type and prints them in markdown table format
func printGroupedCommitsInTableForMarkdown(commits []*object.Commit, config *Config) {
	// Parse and filter commits
	var commitInfos []*CommitInfo
	var breakingChanges []*CommitInfo

	for _, commit := range commits {
		if shouldIgnoreCommit(commit.Message, config.Ignore) {
			continue
		}

		info := parseCommit(commit)
		commitInfos = append(commitInfos, info)

		if info.IsBreaking {
			breakingChanges = append(breakingChanges, info)
		}
	}

	// If all commits were filtered out, show message
	if len(commitInfos) == 0 {
		fmt.Println("_No commits (all filtered)_")
		fmt.Println()
		return
	}

	// Print breaking changes first if any exist
	if len(breakingChanges) > 0 {
		breakingTitle := config.CommitGroups.TitleMaps["breaking"]
		if breakingTitle == "" {
			breakingTitle = "Breaking Changes"
		}

		fmt.Printf("### ⚠️  %s\n\n", breakingTitle)
		fmt.Println("| Commit | Scope | Description |")
		fmt.Println("|--------|-------|-------------|")
		for _, info := range breakingChanges {
			message := info.Subject
			if info.BreakingMessage != "" {
				message = info.BreakingMessage
			}

			scope := "-"
			if info.Scope != "" {
				scope = info.Scope
			}

			// Escape pipe characters in message to prevent table breaking
			message = escapePipes(message)

			fmt.Printf("| `%s` | %s | %s |\n",
				info.Commit.Hash.String()[:7],
				scope,
				message)
		}
		fmt.Println()
	}

	// Group non-breaking commits by type
	groups := make(map[string][]*CommitInfo)
	for _, info := range commitInfos {
		groups[info.Type] = append(groups[info.Type], info)
	}

	// Define order of groups based on config
	var groupOrder []string
	for key := range config.CommitGroups.TitleMaps {
		if key == "breaking" {
			continue // Already printed above
		}
		if _, exists := groups[key]; exists {
			groupOrder = append(groupOrder, key)
		}
	}
	// Add "other" category at the end if it exists
	if _, exists := groups["other"]; exists {
		groupOrder = append(groupOrder, "other")
	}

	// Print each group
	for _, groupType := range groupOrder {
		title := config.CommitGroups.TitleMaps[groupType]
		if title == "" {
			title = "Other"
		}

		fmt.Printf("### %s\n\n", title)
		fmt.Println("| Commit | Scope | Description |")
		fmt.Println("|--------|-------|-------------|")
		for _, info := range groups[groupType] {
			message := info.Subject

			scope := "-"
			if info.Scope != "" {
				scope = info.Scope
			}

			// Escape pipe characters in message to prevent table breaking
			message = escapePipes(message)

			fmt.Printf("| `%s` | %s | %s |\n",
				info.Commit.Hash.String()[:7],
				scope,
				message)
		}
		fmt.Println()
	}
}

// escapePipes escapes pipe characters in strings to prevent breaking markdown tables
func escapePipes(s string) string {
	return s
}
