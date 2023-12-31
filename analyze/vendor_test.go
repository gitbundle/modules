// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package analyze

import "testing"

func TestIsVendor(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"cache/", true},
		{"random/cache/", true},
		{"cache", false},
		{"dependencies/", true},
		{"Dependencies/", true},
		{"dependency/", false},
		{"dist/", true},
		{"dist", false},
		{"random/dist/", true},
		{"random/dist", false},
		{"deps/", true},
		{"configure", true},
		{"a/configure", true},
		{"config.guess", true},
		{"config.guess/", false},
		{".vscode/", true},
		{"doc/_build/", true},
		{"a/docs/_build/", true},
		{"a/dasdocs/_build-vsdoc.js", true},
		{"a/dasdocs/_build-vsdoc.j", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := IsVendor(tt.path); got != tt.want {
				t.Errorf("IsVendor() = %v, want %v", got, tt.want)
			}
		})
	}
}
