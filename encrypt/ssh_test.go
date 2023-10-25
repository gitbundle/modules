// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encrypt

import (
	"bytes"
	"testing"
)

func Test_SshEncDec(t *testing.T) {
	pub, priv, err := GenerateSshKeyPairs()
	if err != nil {
		panic(err)
	}
	// pub = strings.TrimPrefix(pub, "ssh-rsa ")

	for _, data := range [][]byte{
		[]byte("hello test"),
		[]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	} {
		encrypted, err := EncryptWithSshKey(data, []byte(pub))
		if err != nil {
			panic(err)
		}

		data2, err := DecryptWithSshKey(encrypted, priv)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(data, data2) {
			panic("missmatch")
		}
	}
}
