// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package npm

// TagProperty is the name of the property for tag management
const TagProperty = "npm.tag"

// Metadata represents the metadata of a npm package
type Metadata struct {
	Scope                   string            `json:"scope,omitempty"`
	Name                    string            `json:"name,omitempty"`
	Description             string            `json:"description,omitempty"`
	Author                  string            `json:"author,omitempty"`
	License                 string            `json:"license,omitempty"`
	ProjectURL              string            `json:"project_url,omitempty"`
	Keywords                []string          `json:"keywords,omitempty"`
	Dependencies            map[string]string `json:"dependencies,omitempty"`
	DevelopmentDependencies map[string]string `json:"development_dependencies,omitempty"`
	PeerDependencies        map[string]string `json:"peer_dependencies,omitempty"`
	OptionalDependencies    map[string]string `json:"optional_dependencies,omitempty"`
	Readme                  string            `json:"readme,omitempty"`
}
