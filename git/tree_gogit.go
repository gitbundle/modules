// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2015 The Gogs Authors. All rights reserved.

//go:build gogit

package git

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Tree represents a flat directory listing.
type Tree struct {
	ID         SHA1
	ResolvedID SHA1
	repo       *Repository

	gogitTree *object.Tree

	// parent tree
	ptree *Tree
}

func (t *Tree) loadTreeObject() error {
	gogitTree, err := t.repo.gogitRepo.TreeObject(t.ID)
	if err != nil {
		return err
	}

	t.gogitTree = gogitTree
	return nil
}

// ListEntries returns all entries of current tree.
func (t *Tree) ListEntries() (Entries, error) {
	if t.gogitTree == nil {
		err := t.loadTreeObject()
		if err != nil {
			return nil, err
		}
	}

	entries := make([]*TreeEntry, len(t.gogitTree.Entries))
	for i, entry := range t.gogitTree.Entries {
		entries[i] = &TreeEntry{
			ID:             entry.Hash,
			gogitTreeEntry: &t.gogitTree.Entries[i],
			ptree:          t,
		}
	}

	return entries, nil
}

// ListEntriesRecursive returns all entries of current tree recursively including all subtrees
func (t *Tree) ListEntriesRecursive() (Entries, error) {
	if t.gogitTree == nil {
		err := t.loadTreeObject()
		if err != nil {
			return nil, err
		}
	}

	var entries []*TreeEntry
	seen := map[plumbing.Hash]bool{}
	walker := object.NewTreeWalker(t.gogitTree, true, seen)
	for {
		fullName, entry, err := walker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if seen[entry.Hash] {
			continue
		}

		convertedEntry := &TreeEntry{
			ID:             entry.Hash,
			gogitTreeEntry: &entry,
			ptree:          t,
			fullName:       fullName,
		}
		entries = append(entries, convertedEntry)
	}

	return entries, nil
}
