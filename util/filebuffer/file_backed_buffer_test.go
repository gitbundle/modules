// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filebuffer

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileBackedBuffer(t *testing.T) {
	cases := []struct {
		MaxMemorySize int
		Data          string
	}{
		{5, "test"},
		{5, "testtest"},
	}

	for _, c := range cases {
		buf, err := CreateFromReader(strings.NewReader(c.Data), c.MaxMemorySize)
		assert.NoError(t, err)

		assert.EqualValues(t, len(c.Data), buf.Size())

		data, err := io.ReadAll(buf)
		assert.NoError(t, err)
		assert.Equal(t, c.Data, string(data))

		assert.NoError(t, buf.Close())
	}
}
