// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

var entry = &Entry{
	Data: map[string]interface{}{},
	pool: &sync.Pool{},
}

type Fields map[string]interface{}

type Entry struct {
	Data map[string]interface{}
	pool *sync.Pool
}

func NewEntry() *Entry {
	return entry
}

func (e *Entry) newEntry() *Entry {
	entry, ok := e.pool.Get().(*Entry)
	if ok {
		return entry
	}
	return &Entry{Data: map[string]interface{}{}, pool: e.pool}
}

func (e *Entry) releaseEntry(entry *Entry) {
	e.Data = map[string]interface{}{}
	e.pool.Put(entry)
}

func (e *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(e.Data)+len(fields))
	for k, v := range e.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Data: data, pool: e.pool}
}

func (e *Entry) WithField(key string, value interface{}) *Entry {
	entry := e.newEntry()
	defer e.releaseEntry(entry)
	return entry.WithFields(Fields{key: value})
}

func (e *Entry) WithError(err error) *Entry {
	entry := e.newEntry()
	defer e.releaseEntry(entry)
	return entry.WithField("error", err)
}

func (e *Entry) log(level Level, format string, v ...interface{}) {
	elems := make([]string, 0, 8)
	for k, v := range e.Data {
		elems = append(elems, fmt.Sprintf("%s=%v", k, v))
	}
	sort.Strings(elems)
	v = append(v, strings.Join(elems, "\t"))
	Log(2, level, format+" %s", v...)
}

func (e *Entry) Trace(format string, v ...interface{}) {
	e.log(TRACE, format, v...)
}

func (e *Entry) Debug(format string, v ...interface{}) {
	e.log(DEBUG, format, v...)
}

func (e *Entry) Info(format string, v ...interface{}) {
	e.log(INFO, format, v...)
}

func (e *Entry) Warn(format string, v ...interface{}) {
	e.log(WARN, format, v...)
}

func (e *Entry) Error(format string, v ...interface{}) {
	e.log(ERROR, format, v...)
}

func (e *Entry) Critical(format string, v ...interface{}) {
	e.log(CRITICAL, format, v...)
}

func (e *Entry) Fatal(format string, v ...interface{}) {
	e.log(FATAL, format, v...)
	Close()
	os.Exit(1)
}

func WithField(key string, value interface{}) *Entry {
	return entry.WithField(key, value)
}

func WithFields(fields Fields) *Entry {
	return entry.WithFields(fields)
}

func WithError(err error) *Entry {
	return entry.WithError(err)
}
