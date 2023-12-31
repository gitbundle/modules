// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conan

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	name             = "ConanPackage"
	version          = "1.2"
	license          = "MIT"
	author           = "GitBundle <info@gitbundle.com>"
	homepage         = "https://gitbundle.com/"
	url              = "https://gitbundle.com/"
	description      = "Description of ConanPackage"
	topic1           = "gitbundle"
	topic2           = "conan"
	contentConanfile = `from conans import ConanFile, CMake, tools

class ConanPackageConan(ConanFile):
    name = "` + name + `"
    version = "` + version + `"
    license = "` + license + `"
    author = "` + author + `"
    homepage = "` + homepage + `"
    url = "` + url + `"
    description = "` + description + `"
    topics = ("` + topic1 + `", "` + topic2 + `")
    settings = "os", "compiler", "build_type", "arch"
    options = {"shared": [True, False], "fPIC": [True, False]}
    default_options = {"shared": False, "fPIC": True}
    generators = "cmake"
`
)

func TestParseConanfile(t *testing.T) {
	metadata, err := ParseConanfile(strings.NewReader(contentConanfile))
	assert.Nil(t, err)
	assert.Equal(t, license, metadata.License)
	assert.Equal(t, author, metadata.Author)
	assert.Equal(t, homepage, metadata.ProjectURL)
	assert.Equal(t, url, metadata.RepositoryURL)
	assert.Equal(t, description, metadata.Description)
	assert.Equal(t, []string{topic1, topic2}, metadata.Keywords)
}
