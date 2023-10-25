// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timeutil

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gitbundle/modules/log"
)

var (
	executablModTime     = time.Now()
	executablModTimeOnce sync.Once
)

// GetExecutableModTime get executable file modified time of current process.
func GetExecutableModTime() time.Time {
	executablModTimeOnce.Do(func() {
		exePath, err := os.Executable()
		if err != nil {
			log.Error("os.Executable: %v", err)
			return
		}

		exePath, err = filepath.Abs(exePath)
		if err != nil {
			log.Error("filepath.Abs: %v", err)
			return
		}

		exePath, err = filepath.EvalSymlinks(exePath)
		if err != nil {
			log.Error("filepath.EvalSymlinks: %v", err)
			return
		}

		st, err := os.Stat(exePath)
		if err != nil {
			log.Error("os.Stat: %v", err)
			return
		}

		executablModTime = st.ModTime()
	})
	return executablModTime
}
