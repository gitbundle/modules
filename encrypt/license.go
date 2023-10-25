// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encrypt

import (
	"encoding/base64"
	"math/rand"
	"time"
)

func EncryptLicense(plainText, sourceKey []byte, shiftKey byte, shiftArray []int64) (string, error) {
	source := transferBytes(sourceKey, shiftArray)
	key := transformBytes(source, shiftKey)
	cipherText, err := AesGcmEncrypt(key, plainText)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptLicense(cipherString string, sourceKey []byte, shiftKey byte, shiftArray []int64) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherString)
	if err != nil {
		return nil, err
	}
	source := transferBytes(sourceKey, shiftArray)
	key := transformBytes(source, shiftKey)
	return AesGcmDecrypt(key, cipherText)
}

func transformBytes(source []byte, shift byte) []byte {
	result := make([]byte, len(source))

	for i, r := range source {
		ascii := (r + shift) % 128
		result[i] = ascii + 128
	}

	return result
}

func transferBytes(source []byte, indexArray []int64) []byte {
	result := make([]byte, len(indexArray))
	sourceLen := len(source)

	for i, index := range indexArray {
		result[i] = source[int(index)%sourceLen]
	}

	return result
}

func generateRandomArray(length int) []int {
	rand.Seed(time.Now().UnixNano())
	randomArray := make([]int, length)

	for i := 0; i < length; i++ {
		randomArray[i] = rand.Int()
	}

	return randomArray
}
