package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
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

// getCommitsBetween returns all commits between fromHash (exclusive) and toHash (inclusive)
func getCommitsBetween(repo *git.Repository, fromHash *plumbing.Hash, toHash plumbing.Hash) ([]*object.Commit, error) {
	var commits []*object.Commit

	// Start from the newer commit
	commitIter, err := repo.Log(&git.LogOptions{From: toHash})
	if err != nil {
		return nil, err
	}
	defer commitIter.Close()

	if fromHash != nil {
		// Collect commits until we reach the older commit
		err = commitIter.ForEach(func(c *object.Commit) error {
			// Stop when we reach the older commit
			if c.Hash == *fromHash {
				return io.EOF
			}
			commits = append(commits, c)
			return nil
		})

		if err != nil && err != io.EOF {
			return nil, err
		}

		return commits, nil
	}

	if fromHash == nil {
		// Collect all commits (for the oldest tag)
		err = commitIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})

		if err != nil && err != io.EOF {
			return nil, err
		}

		return commits, nil
	}

	return nil, fmt.Errorf("unexpected error in getCommitsBetween")
}
