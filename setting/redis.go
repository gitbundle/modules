// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import "github.com/gitbundle/modules/log"

var (
	Redis = struct {
		Connection string
		Debug      bool
		Disable    bool
	}{}
)

func newRedisService() {
	if err := Cfg.Section("redis").MapTo(&Redis); err != nil {
		log.Fatal("Failed to map redis settings: %v", err)
	}

	if Redis.Connection == "" {
		log.Fatal("Empty Connection for redis setting")
	}
}
