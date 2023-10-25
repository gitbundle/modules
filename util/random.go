// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"math/rand"
	"time"
)

func GenerateRandomString() string {
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, randomLength)
	sourceLen := len(sourceString)

	for i := 0; i < randomLength; i++ {
		randomIndex := rand.Intn(sourceLen)
		result[i] = sourceString[randomIndex]
	}

	return string(result)
}

func GenerateRandomStringWithLength(length int) string {
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)
	sourceLen := len(sourceString)

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(sourceLen)
		result[i] = sourceString[randomIndex]
	}

	return string(result)
}

const (
	sourceString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomLength = 64
)
