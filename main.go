package main

import (
	"context"
	"log"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/urfave/cli/v3"
)

const CurrentFlag = "current"

func main() {
	cmd := &cli.Command{
		Name:  "changelog",
		Usage: "Git Changelog Generator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  CurrentFlag,
				Value: "",
				Usage: "mark version tag to current release",
			},
		},
		Action: action,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx context.Context, cmd *cli.Command) error {
	// Load changelog config
	config, err := NewConfig(".changelog.yml")
	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(config.GitPath)
	if err != nil {
		return err
	}

	// Get all tag references
	tagRefs, err := repo.Tags()
	if err != nil {
		return err
	}

	// Collect all tags with their commit times
	tags, err := LoadTags(repo, tagRefs)
	if err != nil {
		return err
	}

	p := NewMarkdownPrinter(config, repo)
	p.MapData(tags)
	p.Print(cmd.String(CurrentFlag))

	return nil
}
