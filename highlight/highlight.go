// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2015 The Gogs Authors. All rights reserved.

package highlight

import (
	"bufio"
	"bytes"
	"fmt"
	gohtml "html"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gitbundle/modules/analyze"
	"github.com/gitbundle/modules/log"
	"github.com/gitbundle/modules/setting"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	lru "github.com/hashicorp/golang-lru"
)

// don't index files larger than this many bytes for performance purposes
const sizeLimit = 1024 * 1024

var (
	// For custom user mapping
	highlightMapping = map[string]string{}

	once sync.Once

	cache *lru.TwoQueueCache
)

// NewContext loads custom highlight map from local config
func NewContext() {
	once.Do(func() {
		keys := setting.Cfg.Section("highlight.mapping").Keys()
		for i := range keys {
			highlightMapping[keys[i].Name()] = keys[i].Value()
		}

		// The size 512 is simply a conservative rule of thumb
		c, err := lru.New2Q(512)
		if err != nil {
			panic(fmt.Sprintf("failed to initialize LRU cache for highlighter: %s", err))
		}
		cache = c
	})
}

// Code returns a HTML version of code string with chroma syntax highlighting classes
func Code(fileName, language, code string) string {
	NewContext()

	// diff view newline will be passed as empty, change to literal '\n' so it can be copied
	// preserve literal newline in blame view
	if code == "" || code == "\n" {
		return "\n"
	}

	if len(code) > sizeLimit {
		return code
	}

	var lexer chroma.Lexer

	if len(language) > 0 {
		lexer = lexers.Get(language)

		if lexer == nil {
			// Attempt stripping off the '?'
			if idx := strings.IndexByte(language, '?'); idx > 0 {
				lexer = lexers.Get(language[:idx])
			}
		}
	}

	if lexer == nil {
		if val, ok := highlightMapping[filepath.Ext(fileName)]; ok {
			// use mapped value to find lexer
			lexer = lexers.Get(val)
		}
	}

	if lexer == nil {
		if l, ok := cache.Get(fileName); ok {
			lexer = l.(chroma.Lexer)
		}
	}

	if lexer == nil {
		lexer = lexers.Match(fileName)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		cache.Add(fileName, lexer)
	}
	return CodeFromLexer(lexer, code)
}

type nopPreWrapper struct{}

func (nopPreWrapper) Start(code bool, styleAttr string) string { return "" }
func (nopPreWrapper) End(code bool) string                     { return "" }

// CodeFromLexer returns a HTML version of code string with chroma syntax highlighting classes
func CodeFromLexer(lexer chroma.Lexer, code string) string {
	formatter := html.New(html.WithClasses(true),
		html.WithLineNumbers(false),
		html.PreventSurroundingPre(true),
	)

	htmlbuf := bytes.Buffer{}
	htmlw := bufio.NewWriter(&htmlbuf)

	iterator, err := lexer.Tokenise(nil, string(code))
	if err != nil {
		log.Error("Can't tokenize code: %v", err)
		return code
	}
	// style not used for live site but need to pass something
	err = formatter.Format(htmlw, styles.GitHub, iterator)
	if err != nil {
		log.Error("Can't format code: %v", err)
		return code
	}

	_ = htmlw.Flush()
	// Chroma will add newlines for certain lexers in order to highlight them properly
	// Once highlighted, strip them here, so they don't cause copy/paste trouble in HTML output
	return strings.TrimSuffix(htmlbuf.String(), "\n")
}

// File returns a slice of chroma syntax highlighted lines of code
func File(numLines int, fileName, language string, code []byte) []string {
	NewContext()

	if len(code) > sizeLimit {
		return plainText(string(code), numLines)
	}
	formatter := html.New(html.WithClasses(true),
		html.WithLineNumbers(false),
		html.WithPreWrapper(nopPreWrapper{}),
	)

	if formatter == nil {
		log.Error("Couldn't create chroma formatter")
		return plainText(string(code), numLines)
	}

	htmlbuf := bytes.Buffer{}
	htmlw := bufio.NewWriter(&htmlbuf)

	var lexer chroma.Lexer

	// provided language overrides everything
	if len(language) > 0 {
		lexer = lexers.Get(language)
	}

	if lexer == nil {
		if val, ok := highlightMapping[filepath.Ext(fileName)]; ok {
			lexer = lexers.Get(val)
		}
	}

	if lexer == nil {
		language := analyze.GetCodeLanguage(fileName, code)

		lexer = lexers.Get(language)
		if lexer == nil {
			lexer = lexers.Match(fileName)
			if lexer == nil {
				lexer = lexers.Fallback
			}
		}
	}

	iterator, err := lexer.Tokenise(nil, string(code))
	if err != nil {
		log.Error("Can't tokenize code: %v", err)
		return plainText(string(code), numLines)
	}

	err = formatter.Format(htmlw, styles.GitHub, iterator)
	if err != nil {
		log.Error("Can't format code: %v", err)
		return plainText(string(code), numLines)
	}

	_ = htmlw.Flush()
	finalNewLine := false
	if len(code) > 0 {
		finalNewLine = code[len(code)-1] == '\n'
	}

	m := strings.SplitN(htmlbuf.String(), `</span></span><span class="line"><span class="cl">`, numLines)
	if len(m) > 0 {
		m[0] = m[0][len(`<span class="line"><span class="cl">`):]
		last := m[len(m)-1]
		m[len(m)-1] = last[:len(last)-len(`</span></span>`)]
	}

	if finalNewLine {
		m = append(m, "<span class=\"w\">\n</span>")
	}

	return m
}

// return unhiglighted map
func plainText(code string, numLines int) []string {
	m := strings.SplitN(code, "\n", numLines)

	for i, content := range m {
		// need to keep lines that are only \n so copy/paste works properly in browser
		if content == "" {
			content = "\n"
		}
		m[i] = gohtml.EscapeString(content)
	}
	return m
}
