// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2015 The Gogs Authors. All rights reserved.

package git

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// EmptySHA defines empty git SHA
const EmptySHA = "0000000000000000000000000000000000000000"

// EmptyTreeSHA is the SHA of an empty tree
const EmptyTreeSHA = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

// SHAPattern can be used to determine if a string is an valid sha
var SHAPattern = regexp.MustCompile(`^[0-9a-f]{4,40}$`)

// MustID always creates a new SHA1 from a [20]byte array with no validation of input.
func MustID(b []byte) SHA1 {
	var id SHA1
	copy(id[:], b)
	return id
}

// NewID creates a new SHA1 from a [20]byte array.
func NewID(b []byte) (SHA1, error) {
	if len(b) != 20 {
		return SHA1{}, fmt.Errorf("Length must be 20: %v", b)
	}
	return MustID(b), nil
}

// MustIDFromString always creates a new sha from a ID with no validation of input.
func MustIDFromString(s string) SHA1 {
	b, _ := hex.DecodeString(s)
	return MustID(b)
}

// NewIDFromString creates a new SHA1 from a ID string of length 40.
func NewIDFromString(s string) (SHA1, error) {
	var id SHA1
	s = strings.TrimSpace(s)
	if len(s) != 40 {
		return id, fmt.Errorf("Length must be 40: %s", s)
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return id, err
	}
	return NewID(b)
}
