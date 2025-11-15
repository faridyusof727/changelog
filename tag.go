package main

import (
	"fmt"
	"sort"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
)

// TagInfo represents information about a Git tag.
type TagInfo struct {
	Name   string
	Hash   plumbing.Hash
	Time   int64
	Commit *object.Commit
}

// LoadTags loads tag information from the given repository and tag references.
// It iterates over each tag reference, resolves it to a commit (handling annotated tags),
// and collects the tag name, hash, commit time, and commit object.
func LoadTags(repo *git.Repository, tagRefs storer.ReferenceIter) ([]*TagInfo, error) {
	var tagInfos []*TagInfo
	err := tagRefs.ForEach(func(ref *plumbing.Reference) error {
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

		tagInfos = append(tagInfos, &TagInfo{
			Name:   ref.Name().Short(),
			Hash:   ref.Hash(),
			Time:   commit.Committer.When.Unix(),
			Commit: commit,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate tags: %w", err)
	}

	sort.Slice(tagInfos, func(i, j int) bool {
		return tagInfos[i].Time > tagInfos[j].Time
	})

	return tagInfos, nil
}
