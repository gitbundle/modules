// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import "github.com/yuin/goldmark/util"

func isSnakeCaseUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

func isSnakeCaseLowerOrNumber(c byte) bool {
	return 'a' <= c && c <= 'z' || '0' <= c && c <= '9'
}

// ToSnakeCase convert the input string to snake_case format.
//
// Some samples.
//
//	"FirstName"  => "first_name"
//	"HTTPServer" => "http_server"
//	"NoHTTPS"    => "no_https"
//	"GO_PATH"    => "go_path"
//	"GO PATH"    => "go_path"      // space is converted to underscore.
//	"GO-PATH"    => "go_path"      // hyphen is converted to underscore.
func ToSnakeCase(input string) string {
	if len(input) == 0 {
		return ""
	}

	var res []byte
	if len(input) == 1 {
		c := input[0]
		if isSnakeCaseUpper(c) {
			res = []byte{c + 'a' - 'A'}
		} else if isSnakeCaseLowerOrNumber(c) {
			res = []byte{c}
		} else {
			res = []byte{'_'}
		}
	} else {
		res = make([]byte, 0, len(input)*4/3)
		pos := 0
		needSep := false
		for pos < len(input) {
			c := input[pos]
			if c >= 0x80 {
				res = append(res, c)
				pos++
				continue
			}
			isUpper := isSnakeCaseUpper(c)
			if isUpper || isSnakeCaseLowerOrNumber(c) {
				end := pos + 1
				if isUpper {
					// skip the following upper letters
					for end < len(input) && isSnakeCaseUpper(input[end]) {
						end++
					}
					if end-pos > 1 && end < len(input) && isSnakeCaseLowerOrNumber(input[end]) {
						end--
					}
				}
				// skip the following lower or number letters
				for end < len(input) && (isSnakeCaseLowerOrNumber(input[end]) || input[end] >= 0x80) {
					end++
				}
				if needSep {
					res = append(res, '_')
				}
				res = append(res, input[pos:end]...)
				pos = end
				needSep = true
			} else {
				res = append(res, '_')
				pos++
				needSep = false
			}
		}
		for i := 0; i < len(res); i++ {
			if isSnakeCaseUpper(res[i]) {
				res[i] += 'a' - 'A'
			}
		}
	}
	return util.BytesToReadOnlyString(res)
}
