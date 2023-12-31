// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/gitbundle/modules/setting"
)

var externalTrackerRegex = regexp.MustCompile(`({?)(?:user|repo|index)+?(}?)`)

func isLoopbackIP(ip string) bool {
	return net.ParseIP(ip).IsLoopback()
}

// IsValidURL checks if URL is valid
func IsValidURL(uri string) bool {
	if u, err := url.ParseRequestURI(uri); err != nil ||
		(u.Scheme != "http" && u.Scheme != "https") ||
		!validPort(portOnly(u.Host)) {
		return false
	}

	return true
}

// IsValidSiteURL checks if URL is valid
func IsValidSiteURL(uri string) bool {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	if !validPort(portOnly(u.Host)) {
		return false
	}

	for _, scheme := range setting.Service.ValidSiteURLSchemes {
		if scheme == u.Scheme {
			return true
		}
	}
	return false
}

// IsAPIURL checks if URL is current GitBundle instance API URL
func IsAPIURL(uri string) bool {
	return strings.HasPrefix(strings.ToLower(uri), strings.ToLower(setting.AppURL+"api"))
}

// IsValidExternalURL checks if URL is valid external URL
func IsValidExternalURL(uri string) bool {
	if !IsValidURL(uri) || IsAPIURL(uri) {
		return false
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	// Currently check only if not loopback IP is provided to keep compatibility
	if isLoopbackIP(u.Hostname()) || strings.ToLower(u.Hostname()) == "localhost" {
		return false
	}

	// TODO: Later it should be added to allow local network IP addresses
	//       only if allowed by special setting

	return true
}

// IsValidExternalTrackerURLFormat checks if URL matches required syntax for external trackers
func IsValidExternalTrackerURLFormat(uri string) bool {
	if !IsValidExternalURL(uri) {
		return false
	}

	// check for typoed variables like /{index/ or /[repo}
	for _, match := range externalTrackerRegex.FindAllStringSubmatch(uri, -1) {
		if (match[1] == "{" || match[2] == "}") && (match[1] != "{" || match[2] != "}") {
			return false
		}
	}

	return true
}
