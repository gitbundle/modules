// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildLocalPath(t *testing.T) {
	kases := []struct {
		localDir string
		path     string
		expected string
	}{
		{
			"a",
			"0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
			"a/0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
		},
		{
			"a",
			"../0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
			"a/0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
		},
		{
			"a",
			"0\\a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
			"a/0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
		},
		{
			"b",
			"a/../0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
			"b/0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
		},
		{
			"b",
			"a\\..\\0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
			"b/0/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14",
		},
	}

	for _, k := range kases {
		t.Run(k.path, func(t *testing.T) {
			l := LocalStorage{dir: k.localDir}

			assert.EqualValues(t, k.expected, l.buildLocalPath(k.path))
		})
	}
}
