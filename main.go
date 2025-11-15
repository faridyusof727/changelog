package main

import (
	"fmt"
	"io"
	"sort"

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
