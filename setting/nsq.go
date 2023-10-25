// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import "github.com/gitbundle/modules/log"

var (
	Nsq = struct {
		ClusterTcpAddr string
		AuthSecret     string
		Debug          bool
		Disable        bool
	}{}
)

func newNsqService() {
	if err := Cfg.Section("nsq").MapTo(&Nsq); err != nil {
		log.Fatal("Failed to map nsq settings: %v", err)
	}

	if Nsq.ClusterTcpAddr == "" {
		log.Fatal("Empty ClusterTcpAddr for nsq setting")
	}
}
