package main

import (
	"github.com/go-git/go-git/v6"
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
	tags, err := LoadTags(repo, tagRefs)
	if err != nil {
		panic(err)
	}

	p := NewMarkdownPrinter(config, repo)
	p.MapData(tags)
	p.Print()
}
