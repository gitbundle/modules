// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"testing"
)

func Test_GenerateHashString(t *testing.T) {
	salt, err := CryptoRandomString(32)
	if err != nil {
		t.Error(err)
	}
	token := GenerateRandomString()
	hashedString := GenerateHashString(token, salt)
	fmt.Println(hashedString, len(hashedString))
	if len(hashedString) > 255 {
		t.Error("too long hash string")
	}
}
