// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package regexplru

import (
	"regexp"

	"github.com/gitbundle/modules/log"

	lru "github.com/hashicorp/golang-lru"
)

var lruCache *lru.Cache

func init() {
	var err error
	lruCache, err = lru.New(1000)
	if err != nil {
		log.Fatal("failed to new LRU cache, err: %v", err)
	}
}

// GetCompiled works like regexp.Compile, the compiled expr or error is stored in LRU cache
func GetCompiled(expr string) (r *regexp.Regexp, err error) {
	v, ok := lruCache.Get(expr)
	if !ok {
		r, err = regexp.Compile(expr)
		if err != nil {
			lruCache.Add(expr, err)
			return nil, err
		}
		lruCache.Add(expr, r)
	} else {
		r, ok = v.(*regexp.Regexp)
		if !ok {
			if err, ok = v.(error); ok {
				return nil, err
			}
			panic("impossible")
		}
	}
	return r, nil
}
