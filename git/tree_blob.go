// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2015 The Gogs Authors. All rights reserved.

package git

// GetBlobByPath get the blob object according the path
func (t *Tree) GetBlobByPath(relpath string) (*Blob, error) {
	entry, err := t.GetTreeEntryByPath(relpath)
	if err != nil {
		return nil, err
	}

	if !entry.IsDir() && !entry.IsSubModule() {
		return entry.Blob(), nil
	}

	return nil, ErrNotExist{"", relpath}
}
