// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"context"

	giturl "github.com/gitbundle/modules/git/url"
)

// GetRemoteAddress returns remote url of git repository in the repoPath with special remote name
func GetRemoteAddress(ctx context.Context, repoPath, remoteName string) (string, error) {
	var cmd *Command
	if CheckGitVersionAtLeast("2.7") == nil {
		cmd = NewCommand(ctx, "remote", "get-url", remoteName)
	} else {
		cmd = NewCommand(ctx, "config", "--get", "remote."+remoteName+".url")
	}

	result, _, err := cmd.RunStdString(&RunOpts{Dir: repoPath})
	if err != nil {
		return "", err
	}

	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return result, nil
}

// GetRemoteURL returns the url of a specific remote of the repository.
func GetRemoteURL(ctx context.Context, repoPath, remoteName string) (*giturl.GitURL, error) {
	addr, err := GetRemoteAddress(ctx, repoPath, remoteName)
	if err != nil {
		return nil, err
	}
	return giturl.Parse(addr)
}
