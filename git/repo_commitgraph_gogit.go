// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2019 The GitBundle Authors.
// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build gogit

package git

import (
	"os"
	"path"

	gitealog "github.com/gitbundle/modules/log"

	"github.com/go-git/go-git/v5/plumbing/format/commitgraph"
	cgobject "github.com/go-git/go-git/v5/plumbing/object/commitgraph"
)

// CommitNodeIndex returns the index for walking commit graph
func (r *Repository) CommitNodeIndex() (cgobject.CommitNodeIndex, *os.File) {
	indexPath := path.Join(r.Path, "objects", "info", "commit-graph")

	file, err := os.Open(indexPath)
	if err == nil {
		var index commitgraph.Index
		index, err = commitgraph.OpenFileIndex(file)
		if err == nil {
			return cgobject.NewGraphCommitNodeIndex(index, r.gogitRepo.Storer), file
		}
	}

	if !os.IsNotExist(err) {
		gitealog.Warn("Unable to read commit-graph for %s: %v", r.Path, err)
	}

	return cgobject.NewObjectCommitNodeIndex(r.gogitRepo.Storer), nil
}
