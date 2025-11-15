package main

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func main() {
	// Load changelog config
	config, err := NewConfig(".changelog.yml")
	if err != nil {
		panic(err)
	}

	repo, err := git.PlainOpen(config.GitPath)
	if err != nil {
		panic(err)
	}

	// Get all tag references
	tagRefs, err := repo.Tags()
	if err != nil {
		panic(err)
	}

	// Collect all tags with their commit times
	var tags []TagInfo
	err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
		// Get the commit that the tag points to
		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			// If it's an annotated tag, resolve it
			tag, err := repo.TagObject(ref.Hash())
			if err != nil {
				return err
			}
			commit, err = tag.Commit()
			if err != nil {
				return err
			}
		}

		tags = append(tags, TagInfo{
			Name:   ref.Name().Short(),
			Hash:   ref.Hash(),
			Time:   commit.Committer.When.Unix(),
			Commit: commit,
		})
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Sort by time (newest first)
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Time > tags[j].Time
	})

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

func getFirstLine(message string) string {
	message = strings.TrimSpace(message)
	if idx := strings.Index(message, "\n"); idx != -1 {
		return strings.TrimSpace(message[:idx])
	}
	return message
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

		// Collect commits until we reach the older commit
		err = commitIter.ForEach(func(c *object.Commit) error {
			commits = append(commits, c)
			return nil
		})

		if err != nil && err != io.EOF {
			return nil, err
		}

		return commits, nil

	}

	return nil, fmt.Errorf("error")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// extractCommitType extracts the type from a conventional commit message
// e.g., "feat: add new feature" -> "feat"
func extractCommitType(message string) string {
	// Match conventional commit format: type(scope): message or type: message
	re := regexp.MustCompile(`^(\w+)(?:\([^)]+\))?:\s`)
	matches := re.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return "other"
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
	// Filter out ignored commits
	var filteredCommits []*object.Commit
	for _, commit := range commits {
		if !shouldIgnoreCommit(commit.Message, config.Ignore) {
			filteredCommits = append(filteredCommits, commit)
		}
	}

	// If all commits were filtered out, show message
	if len(filteredCommits) == 0 {
		fmt.Println("  No commits (all filtered)")
		return
	}

	// Group commits by type
	groups := make(map[string][]*object.Commit)
	for _, commit := range filteredCommits {
		commitType := extractCommitType(getFirstLine(commit.Message))
		groups[commitType] = append(groups[commitType], commit)
	}

	// Define order of groups based on config
	var groupOrder []string
	for key := range config.CommitGroups.TitleMaps {
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
		for _, commit := range groups[groupType] {
			message := getFirstLine(commit.Message)
			// Remove the type prefix for cleaner output
			message = regexp.MustCompile(`^\w+(?:\([^)]+\))?:\s*`).ReplaceAllString(message, "")

			fmt.Printf("  • %s - %s\n",
				commit.Hash.String()[:7],
				message)
		}
		fmt.Println()
	}
}
