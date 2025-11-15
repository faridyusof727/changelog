package main

import (
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type TagInfo struct {
	Name   string
	Hash   plumbing.Hash
	Time   int64
	Commit *object.Commit
}
