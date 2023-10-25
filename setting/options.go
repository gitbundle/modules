// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2016 The GitBundle Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package setting

import (
	"fmt"
	"os"
	"path"

	"github.com/gitbundle/modules/log"
	"github.com/gitbundle/modules/util"
)

var directories = NewDirectorySet()

// Dir returns all files from static or custom directory.
func Dir(name string) ([]string, error) {
	if directories.Filled(name) {
		return directories.Get(name), nil
	}

	var result []string

	customDir := path.Join(CustomPath, "options", name)

	isDir, err := util.IsDir(customDir)
	if err != nil {
		return []string{}, fmt.Errorf("Unabe to check if custom directory %s is a directory. %v", customDir, err)
	}
	if isDir {
		files, err := util.StatDir(customDir, true)
		if err != nil {
			return []string{}, fmt.Errorf("Failed to read custom directory. %v", err)
		}

		result = append(result, files...)
	}

	staticDir := path.Join(StaticRootPath, "options", name)

	isDir, err = util.IsDir(staticDir)
	if err != nil {
		return []string{}, fmt.Errorf("Unabe to check if static directory %s is a directory. %v", staticDir, err)
	}
	if isDir {
		files, err := util.StatDir(staticDir, true)
		if err != nil {
			return []string{}, fmt.Errorf("Failed to read static directory. %v", err)
		}

		result = append(result, files...)
	}

	return directories.AddAndGet(name, result), nil
}

// Locale reads the content of a specific locale from static or custom path.
func Locale(name string) ([]byte, error) {
	return fileFromDir(path.Join("locale", name))
}

// Readme reads the content of a specific readme from static or custom path.
func Readme(name string) ([]byte, error) {
	return fileFromDir(path.Join("readme", name))
}

// Gitignore reads the content of a specific gitignore from static or custom path.
func Gitignore(name string) ([]byte, error) {
	return fileFromDir(path.Join("gitignore", name))
}

// License reads the content of a specific license from static or custom path.
func License(name string) ([]byte, error) {
	return fileFromDir(path.Join("license", name))
}

// Labels reads the content of a specific labels from static or custom path.
func Labels(name string) ([]byte, error) {
	return fileFromDir(path.Join("label", name))
}

// fileFromDir is a helper to read files from static or custom path.
func fileFromDir(name string) ([]byte, error) {
	customPath := path.Join(CustomPath, "options", name)

	isFile, err := util.IsFile(customPath)
	if err != nil {
		log.Error("Unable to check if %s is a file. Error: %v", customPath, err)
	}
	if isFile {
		return os.ReadFile(customPath)
	}

	staticPath := path.Join(StaticRootPath, "options", name)

	isFile, err = util.IsFile(staticPath)
	if err != nil {
		log.Error("Unable to check if %s is a file. Error: %v", staticPath, err)
	}
	if isFile {
		return os.ReadFile(staticPath)
	}

	return []byte{}, fmt.Errorf("Asset file does not exist: %s", name)
}

type directorySet map[string][]string

func NewDirectorySet() directorySet {
	return make(directorySet)
}

func (s directorySet) Add(key string, value []string) {
	_, ok := s[key]

	if !ok {
		s[key] = make([]string, 0, len(value))
	}

	s[key] = append(s[key], value...)
}

func (s directorySet) Get(key string) []string {
	_, ok := s[key]

	if ok {
		result := []string{}
		seen := map[string]string{}

		for _, val := range s[key] {
			if _, ok := seen[val]; !ok {
				result = append(result, val)
				seen[val] = val
			}
		}

		return result
	}

	return []string{}
}

func (s directorySet) AddAndGet(key string, value []string) []string {
	s.Add(key, value)
	return s.Get(key)
}

func (s directorySet) Filled(key string) bool {
	return len(s[key]) > 0
}
