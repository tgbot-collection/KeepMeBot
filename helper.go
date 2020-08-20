// KeepMeBot - helper
// 2020-08-23 17:36
// Benny <benny.think@gmail.com>

package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var version, hash, _ = getTagAndHash()

func getTagAndHash() (string, string, string) {
	r, err := git.PlainOpen(".")
	if err != nil {
		r, _ = git.PlainOpen(".")
	}
	ref, _ := r.Head()
	hash := ref.Hash().String()[:6]
	commit, _ := r.CommitObject(ref.Hash())

	var tag string
	t, _ := r.TagObjects()

	_ = t.ForEach(func(t *object.Tag) error {
		if tag < t.Name {
			tag = t.Name
		}
		return nil
	})
	if tag == "" {
		tag = "No tags🧐"
	}
	return tag, hash, commit.String()
}
