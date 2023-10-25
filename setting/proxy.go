// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	"net/url"

	"github.com/gitbundle/modules/log"
)

// Proxy settings
var Proxy = struct {
	Enabled       bool
	ProxyURL      string
	ProxyURLFixed *url.URL
	ProxyHosts    []string
}{
	Enabled:    false,
	ProxyURL:   "",
	ProxyHosts: []string{},
}

func newProxyService() {
	sec := Cfg.Section("proxy")
	Proxy.Enabled = sec.Key("PROXY_ENABLED").MustBool(false)
	Proxy.ProxyURL = sec.Key("PROXY_URL").MustString("")
	if Proxy.ProxyURL != "" {
		var err error
		Proxy.ProxyURLFixed, err = url.Parse(Proxy.ProxyURL)
		if err != nil {
			log.Error("Global PROXY_URL is not valid")
			Proxy.ProxyURL = ""
		}
	}
	Proxy.ProxyHosts = sec.Key("PROXY_HOSTS").Strings(",")
}
