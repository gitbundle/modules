// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build !windows

package log

import (
	"os"

	"github.com/mattn/go-isatty"
)

func init() {
	// when running gitbundle as a systemd unit with logging set to console, the output can not be colorized,
	// otherwise it spams the journal / syslog with escape sequences like "#033[0m#033[32mcmd/web.go:102:#033[32m"
	// this file covers non-windows platforms.
	CanColorStdout = isatty.IsTerminal(os.Stdout.Fd())
	CanColorStderr = isatty.IsTerminal(os.Stderr.Fd())
}
