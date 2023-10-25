// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	testContent := []byte("hello")

	tmpDir := os.TempDir()
	now := time.Now()
	srcFile := fmt.Sprintf("%s/copy-test-%d-src.txt", tmpDir, now.UnixMicro())
	dstFile := fmt.Sprintf("%s/copy-test-%d-dst.txt", tmpDir, now.UnixMicro())

	_ = os.Remove(srcFile)
	_ = os.Remove(dstFile)
	defer func() {
		_ = os.Remove(srcFile)
		_ = os.Remove(dstFile)
	}()

	err := os.WriteFile(srcFile, testContent, 0o777)
	assert.NoError(t, err)
	err = CopyFile(srcFile, dstFile)
	assert.NoError(t, err)
	dstContent, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, testContent, dstContent)
}
