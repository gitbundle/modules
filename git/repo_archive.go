// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2015 The Gogs Authors. All rights reserved.

package git

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// ArchiveType archive types
type ArchiveType int

const (
	// ZIP zip archive type
	ZIP ArchiveType = iota + 1
	// TARGZ tar gz archive type
	TARGZ
	// BUNDLE bundle archive type
	BUNDLE
)

// String converts an ArchiveType to string
func (a ArchiveType) String() string {
	switch a {
	case ZIP:
		return "zip"
	case TARGZ:
		return "tar.gz"
	case BUNDLE:
		return "bundle"
	}
	return "unknown"
}

// CreateArchive create archive content to the target path
func (repo *Repository) CreateArchive(ctx context.Context, format ArchiveType, target io.Writer, usePrefix bool, commitID string) error {
	if format.String() == "unknown" {
		return fmt.Errorf("unknown format: %v", format)
	}

	args := []string{
		"archive",
	}
	if usePrefix {
		args = append(args, "--prefix="+filepath.Base(strings.TrimSuffix(repo.Path, ".git"))+"/")
	}

	args = append(args,
		"--format="+format.String(),
		commitID,
	)

	var stderr strings.Builder
	err := NewCommand(ctx, args...).Run(&RunOpts{
		Dir:    repo.Path,
		Stdout: target,
		Stderr: &stderr,
	})
	if err != nil {
		return ConcatenateError(err, stderr.String())
	}
	return nil
}
