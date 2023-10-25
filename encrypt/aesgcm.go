// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Aesgcm provides an encrypter that uses the aesgcm encryption
// algorithm.
type Aesgcm struct {
	block  cipher.Block
	Compat bool
}

// Encrypt encrypts the plaintext using aesgcm.
func (e *Aesgcm) Encrypt(plaintext string) ([]byte, error) {
	gcm, err := cipher.NewGCM(e.block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

// Decrypt decrypts the ciphertext using aesgcm.
func (e *Aesgcm) Decrypt(ciphertext []byte) (string, error) {
	gcm, err := cipher.NewGCM(e.block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		// if the decryption utility is running in compatibility
		// mode, it will return the ciphertext as plain text if
		// decryption fails. This should be used when running the
		// database in mixed-mode, where there is a mix of encrypted
		// and unencrypted content.
		if e.Compat {
			return string(ciphertext), nil
		}
		return "", errors.New("malformed ciphertext")
	}

	plaintext, err := gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
	// if the decryption utility is running in compatibility
	// mode, it will return the ciphertext as plain text if
	// decryption fails. This should be used when running the
	// database in mixed-mode, where there is a mix of encrypted
	// and unencrypted content.
	if err != nil && e.Compat {
		return string(ciphertext), nil
	}
	return string(plaintext), err
}

func AesGcmKey32Encrypt(key, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errKeySize
	}
	return AesGcmEncrypt(key, plaintext)
}

func AesGcmKey32Decrypt(key, ciphertext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errKeySize
	}
	return AesGcmDecrypt(key, ciphertext)
}

// AesGcmEncrypt (from legacy package): encrypts plaintext with the given key using AES in GCM mode. should be replaced.
func AesGcmEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

// AesGcmDecrypt (from legacy package): decrypts ciphertext with the given key using AES in GCM mode. should be replaced.
func AesGcmDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	plainText, err := gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
