package main

import (
	"fmt"

	"github.com/go-git/go-git/v6"
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
