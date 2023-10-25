// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	"time"

	"github.com/gitbundle/modules/log"
)

var (
	Deployment = struct {
		DebugRunningDuration  time.Duration
		VerifyRunningDuration time.Duration
		Disable               bool
	}{}
)

func newDeploymentService() {
	sec := Cfg.Section("deployment")
	if err := sec.MapTo(&Deployment); err != nil {
		log.Fatal("Failed to map nsq settings: %v", err)
	}

	Deployment.DebugRunningDuration = sec.Key("DEBUG_RUNNING_DURATION").MustDuration(7 * 24 * time.Hour)
	Deployment.VerifyRunningDuration = sec.Key("VERIFY_RUNNING_DURATION").MustDuration(10 * time.Minute)
}
