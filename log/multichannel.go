// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// MultiChannelledLogger is default logger in the GitBundle application.
// it can contain several providers and log message into all providers.
type MultiChannelledLogger struct {
	LevelLoggerLogger
	*MultiChannelledLog
	bufferLength int64
	noCaller     bool
}

// newLogger initializes and returns a new logger.
func newLogger(name string, buffer int64, noCaller bool) *MultiChannelledLogger {
	l := &MultiChannelledLogger{
		MultiChannelledLog: NewMultiChannelledLog(name, buffer),
		bufferLength:       buffer,
		noCaller:           noCaller,
	}
	l.LevelLogger = l
	return l
}

// SetLogger sets new logger instance with given logger provider and config.
func (l *MultiChannelledLogger) SetLogger(name, provider, config string) error {
	eventLogger, err := NewChannelledLog(l.ctx, name, provider, config, l.bufferLength)
	if err != nil {
		return fmt.Errorf("Failed to create sublogger (%s): %v", name, err)
	}

	l.MultiChannelledLog.DelLogger(name)

	err = l.MultiChannelledLog.AddLogger(eventLogger)
	if err != nil {
		if IsErrDuplicateName(err) {
			return fmt.Errorf("Duplicate named sublogger %s %v", name, l.MultiChannelledLog.GetEventLoggerNames())
		}
		return fmt.Errorf("Failed to add sublogger (%s): %v", name, err)
	}

	return nil
}

// DelLogger deletes a sublogger from this logger.
func (l *MultiChannelledLogger) DelLogger(name string) (bool, error) {
	return l.MultiChannelledLog.DelLogger(name), nil
}

// Log msg at the provided level with the provided caller defined by skip (0 being the function that calls this function)
func (l *MultiChannelledLogger) Log(skip int, level Level, format string, v ...interface{}) error {
	if l.GetLevel() > level {
		return nil
	}

	var (
		caller   string
		stack    string
		filename string
		line     int
		pc       uintptr
		ok       bool
	)

	if !l.noCaller {
		caller = "?()"
		pc, filename, line, ok = runtime.Caller(skip + 1)
		if ok {
			// Get caller function name.
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				caller = fn.Name() + "()"
			}
		}
		if l.GetStacktraceLevel() <= level {
			stack = Stack(skip + 1)
		}
	}
	msg := format
	if len(v) > 0 {
		msg = ColorSprintf(format, v...)
	}
	labels := getGoroutineLabels()
	if labels != nil {
		pid, ok := labels["pid"]
		if ok {
			msg = "[" + ColorString(FgHiYellow) + pid + ColorString(Reset) + "] " + msg
		}
	}

	return l.SendLog(level, caller, strings.TrimPrefix(filename, prefix), line, msg, stack)
}

// SendLog sends a log event at the provided level with the information given
func (l *MultiChannelledLogger) SendLog(level Level, caller, filename string, line int, msg, stack string) error {
	if l.GetLevel() > level {
		return nil
	}
	event := &Event{
		level:      level,
		caller:     caller,
		filename:   filename,
		line:       line,
		msg:        msg,
		time:       time.Now(),
		stacktrace: stack,
	}
	l.LogEvent(event)
	return nil
}
