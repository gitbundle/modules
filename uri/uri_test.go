// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package uri

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadURI(t *testing.T) {
	p, err := filepath.Abs("./uri.go")
	assert.NoError(t, err)
	f, err := Open("file://" + p)
	assert.NoError(t, err)
	defer f.Close()
}
