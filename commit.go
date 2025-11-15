package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v6/plumbing/object"
)

// CommitInfo holds parsed information about a commit
type CommitInfo struct {
	Commit          *object.Commit
	Type            string
	Scope           string
	Subject         string
	Body            string
	IsBreaking      bool
	BreakingMessage string
}

// getFirstLine extracts the first line (subject) from a commit message
func getFirstLine(message string) string {
	message = strings.TrimSpace(message)
	if idx := strings.Index(message, "\n"); idx != -1 {
		return strings.TrimSpace(message[:idx])
	}
	return message
}

// extractCommitType extracts the type from a conventional commit message
// e.g., "feat: add new feature" -> "feat"
func extractCommitType(message string) string {
	// Match conventional commit format: type(scope): message or type: message
	re := regexp.MustCompile(`^(\w+)(?:\([^)]+\))?!?:\s`)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return "other"
}

// parseCommit parses a commit into a CommitInfo struct
func parseCommit(commit *object.Commit) *CommitInfo {
	message := commit.Message
	subject := getFirstLine(message)

	// Extract type, scope, and check for breaking change indicator
	typeRe := regexp.MustCompile(`^(\w+)(?:\(([^)]+)\))?(!)?:\s*(.*)$`)
	matches := typeRe.FindStringSubmatch(subject)

	info := &CommitInfo{
		Commit:  commit,
		Type:    "other",
		Subject: subject,
	}

	if len(matches) > 1 {
		info.Type = matches[1]
		if len(matches) > 2 {
			info.Scope = matches[2]
		}
		if len(matches) > 3 && matches[3] == "!" {
			info.IsBreaking = true
		}
		if len(matches) > 4 {
			info.Subject = matches[4]
		}
	}

	// Extract body and check for BREAKING CHANGE footer
	parts := strings.SplitN(message, "\n", 2)
	if len(parts) > 1 {
		info.Body = strings.TrimSpace(parts[1])

		// Check for BREAKING CHANGE or BREAKING-CHANGE footer
		breakingRe := regexp.MustCompile(`(?m)^BREAKING[- ]CHANGE:\s*(.+)$`)
		if matches := breakingRe.FindStringSubmatch(info.Body); len(matches) > 1 {
			info.IsBreaking = true
			info.BreakingMessage = strings.TrimSpace(matches[1])
		}
	}

	return info
}

// shouldIgnoreCommit checks if a commit message contains the ignore pattern
func shouldIgnoreCommit(message string, ignorePattern string) bool {
	if ignorePattern == "" {
		return false
	}
	return strings.Contains(message, ignorePattern)
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
