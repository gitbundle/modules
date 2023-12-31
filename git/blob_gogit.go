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
)

// Blob represents a Git object.
type Blob struct {
	ID SHA1

	gogitEncodedObj plumbing.EncodedObject
	name            string
}

// DataAsync gets a ReadCloser for the contents of a blob without reading it all.
// Calling the Close function on the result will discard all unread output.
func (b *Blob) DataAsync() (io.ReadCloser, error) {
	return b.gogitEncodedObj.Reader()
}

// Size returns the uncompressed size of the blob
func (b *Blob) Size() int64 {
	return b.gogitEncodedObj.Size()
}
