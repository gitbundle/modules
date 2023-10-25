// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	golog "log"
	"os"
	"path"

	"github.com/gitbundle/modules/json"
)

// mode can be console, file
func NewSimpleLogger(logBufferLength int64, mode, level, stacktraceLevel, logPath string, noCaller bool) {
	lvl := INFO
	if level != "" {
		lvl = FromString(level)
	}

	stacktraceLvl := NONE
	if stacktraceLevel != "" {
		stacktraceLvl = FromString(stacktraceLevel)
	}

	logConfig := map[string]interface{}{
		"level":           lvl.String(),
		"expression":      "",
		"prefix":          "",
		"flags":           LstdFlags,
		"stacktraceLevel": stacktraceLvl.String(),
	}

	// Generate log configuration.
	switch mode {
	case "console":
		logConfig["stderr"] = true
		logConfig["colorize"] = CanColorStderr
	case "file":
		if err := os.MkdirAll(path.Dir(logPath), os.ModePerm); err != nil {
			panic(err.Error())
		}

		logConfig["filename"] = logPath
		logConfig["rotate"] = true
		logConfig["maxsize"] = 1 << 28
		logConfig["daily"] = true
		logConfig["maxdays"] = 7
		logConfig["compress"] = true
		logConfig["compressionLevel"] = -1
		logConfig["colorize"] = false
	}

	byteConfig, err := json.Marshal(logConfig)
	if err != nil {
		panic(err)
	}

	NewLogger(logBufferLength, mode, mode, string(byteConfig), noCaller)

	// Finally redirect the default golog to here
	golog.SetFlags(0)
	golog.SetPrefix("")
	golog.SetOutput(NewLoggerAsWriter("INFO", GetLogger(DEFAULT)))
}
