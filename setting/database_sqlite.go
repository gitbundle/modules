// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build sqlite

// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	EnableSQLite3 = false
	// SupportedDatabaseTypes = append(SupportedDatabaseTypes, "sqlite3")
}
