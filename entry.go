package log

/*
   Copyright 2019 Bruno Moura <brunotm@gmail.com>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	hex = "0123456789abcdef"
)

// Entry is a structured log entry. A entry is not safe for concurrent use.
// A entry must be logged by calling Log(), and cannot be reused after.
type Entry struct {
	keys  []string
	data  []byte
	l     *Logger
	done  bool
	level Level
}

func (e *Entry) open() {
	e.data = append(e.data, '{')
}

func (e *Entry) close() {
	e.data = append(e.data, '}')
}

func (e *Entry) reset() {
	e.l = nil
	e.done = false
	e.data = e.data[:0]
	e.keys = e.keys[:0]
}

// Log logs the current entry. An entry must not be used after calling Log.
func (e *Entry) Log() {
	if e == nil || e.done {
		return
	}
	e.close()
	e.done = true
	e.l.write(e)
}

// Discard the current entry without logging it.
func (e *Entry) discard() {
	if e == nil || e.done {
		return
	}
	e.l.discard(e)
}

// Level returns the log level of current entry.
func (e *Entry) Level() (level Level) {
	return e.level
}

// Bytes return the current entry bytes. This is intended to be used in hooks
// That will be applied after calling Log().
// The retuned []byte is not a copy and must not be modified directly.
func (e *Entry) Bytes() (data []byte) {
	return e.data
}

// Error adds the given error key/value to the log entry
func (e *Entry) Error(key string, err error) (entry *Entry) {
	e.String(key, err.Error())
	return e
}

// Bool adds the given bool key/value to the log entry
func (e *Entry) Bool(key string, value bool) (entry *Entry) {
	if e == nil || e.done {
		return nil
	}

	if e.hasKey(key) {
		return e
	}

	e.addKey(key)
	e.data = strconv.AppendBool(e.data, value)
	return e
}

// Float adds the given float key/value to the log entry
func (e *Entry) Float(key string, value float64) (entry *Entry) {
	if e == nil || e.done {
		return nil
	}

	if e.hasKey(key) {
		return e
	}

	e.addKey(key)
	e.data = strconv.AppendFloat(e.data, value, 'f', -1, 64)
	return e
}

// Int adds the given int key/value to the log entry
func (e *Entry) Int(key string, value int64) (entry *Entry) {
	if e == nil || e.done {
		return nil
	}

	if e.hasKey(key) {
		return e
	}

	e.addKey(key)
	e.data = strconv.AppendInt(e.data, value, 10)
	return e
}

// Uint adds the given uint key/value to the log entry
func (e *Entry) Uint(key string, value uint64) (entry *Entry) {
	if e == nil || e.done {
		return nil
	}

	if e.hasKey(key) {
		return e
	}

	e.addKey(key)
	e.data = strconv.AppendUint(e.data, value, 10)
	return e
}

// String adds the given string key/value to the log entry
func (e *Entry) String(key string, value string) (entry *Entry) {
	if e == nil || e.done {
		return nil
	}

	if e.hasKey(key) {
		return e
	}

	e.addKey(key)
	e.writeString(value)
	return e
}

// Object creates a nested object within the given log entry.
// Calling Log() on the object Entry is not allowed.
func (e *Entry) Object(key string, fn func(*Entry)) (entry *Entry) {
	var sub *Entry
	if e != nil {
		sub = entryPool.Get().(*Entry)
		sub.open()
	}

	fn(sub)

	if e != nil {
		sub.close()
		e.addKey(key)
		e.data = append(e.data, sub.data...)
		sub.discard()
	}
	return e
}

// based on https://golang.org/src/encoding/json/encode.go:884
func (e *Entry) writeString(s string) {
	e.data = append(e.data, '"')

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '\\' && c != '"' {
			e.data = append(e.data, c)
			continue
		}
		switch c {
		case '"', '\\':
			e.data = append(e.data, '\\', '"')
		case '\n':
			e.data = append(e.data, '\\', '\n')
		case '\f':
			e.data = append(e.data, '\\', '\f')
		case '\b':
			e.data = append(e.data, '\\', '\b')
		case '\r':
			e.data = append(e.data, '\\', '\r')
		case '\t':
			e.data = append(e.data, '\\', '\t')
		default:
			e.data = append(e.data, `\u00`...)
			e.data = append(e.data, hex[c>>4], hex[c&0xF])
		}
		continue
	}

	e.data = append(e.data, '"')
}

func (e *Entry) init(logger *Logger, level Level) {

	t := time.Now()
	e.level = level
	e.l = logger

	e.open()
	e.String("level", level.String())

	if e.l.config.EnableTime {

		switch e.l.config.TimeFormat {
		case Unix:
			e.Int("time", t.Unix())
		case UnixMilli:
			e.Int("time", t.UnixNano()/1000000)
		case UnixNano:
			e.Int("time", t.UnixNano())
		default:
			e.addKey("time")
			e.data = append(e.data, '"')
			e.data = t.AppendFormat(e.data, e.l.config.TimeFormat)
			e.data = append(e.data, '"')
		}

	}

	if e.l.config.EnableCaller {
		_, f, l, ok := runtime.Caller(3 + e.l.config.CallerSkip)

		if ok {
			idx := strings.LastIndexByte(f, '/')
			idx = strings.LastIndexByte(f[:idx], '/')
			e.String("caller", f[idx+1:]+":"+strconv.Itoa(l))
		} else if !ok {
			e.String("caller", "???")
		}

	}
}

func (e *Entry) addKey(key string) {

	if len(e.keys) > 0 {
		e.data = append(e.data, ',')
	}

	e.data = append(e.data, '"')
	e.data = append(e.data, key...)
	e.data = append(e.data, '"', ':')

	e.keys = append(e.keys, key)
}

func (e *Entry) hasKey(key string) bool {
	ks := len(e.keys)
	if ks > 0 {
		for i := 0; i < ks; i++ {
			if key == e.keys[i] {
				return true
			}
		}
	}
	return false
}
