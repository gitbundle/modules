// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"sync"

	logger "github.com/gitbundle/modules/log"
	"github.com/gitbundle/modules/setting"
)

var once sync.Once

// CheckLFSVersion will check lfs version, if not satisfied, then disable it.
func CheckLFSVersion() {
	if setting.LFS.StartServer {
		// Disable LFS client hooks if installed for the current OS user
		// Needs at least git v2.1.2
		if CheckGitVersionAtLeast("2.1.2") != nil {
			setting.LFS.StartServer = false
			logger.Error("LFS server support needs at least Git v2.1.2")
		} else {
			once.Do(func() {
				globalCommandArgs = append(globalCommandArgs, "-c", "filter.lfs.required=",
					"-c", "filter.lfs.smudge=", "-c", "filter.lfs.clean=")
			})
		}
	}
}
