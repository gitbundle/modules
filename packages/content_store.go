// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package packages

import (
	"io"
	"path"

	"github.com/gitbundle/modules/storage"
)

// BlobHash256Key is the key to address a blob content
type BlobHash256Key string

// ContentStore is a wrapper around ObjectStorage
type ContentStore struct {
	store storage.ObjectStorage
}

// NewContentStore creates the default package store
func NewContentStore() *ContentStore {
	contentStore := &ContentStore{storage.Packages}
	return contentStore
}

// Get gets a package blob
func (s *ContentStore) Get(key BlobHash256Key) (storage.Object, error) {
	return s.store.Open(keyToRelativePath(key))
}

// Save stores a package blob
func (s *ContentStore) Save(key BlobHash256Key, r io.Reader, size int64) error {
	_, err := s.store.Save(keyToRelativePath(key), r, size)
	return err
}

// Delete deletes a package blob
func (s *ContentStore) Delete(key BlobHash256Key) error {
	return s.store.Delete(keyToRelativePath(key))
}

// keyToRelativePath converts the sha256 key aabb000000... to aa/bb/aabb000000...
func keyToRelativePath(key BlobHash256Key) string {
	return path.Join(string(key)[0:2], string(key)[2:4], string(key))
}
