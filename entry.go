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
// Repeated keys will be ignored.
type Entry struct {
	enc   *encoder
	l     *Logger
	level Level
}

func (e Entry) reset() {
	e.l = nil
	e.enc.reset()
}

// Log logs the current entry. An entry must not be used after calling Log.
func (e Entry) Log() {
	if e.enc != nil {
		e.enc.closeObject()
		e.l.write(e)
	}
}

// Discard the current entry without logging it.
func (e Entry) discard() {
	if e.enc != nil {
		e.l.discard(e)
	}
}

// Level returns the log level of current entry.
func (e Entry) Level() (level Level) {
	return e.level
}

// Bytes return the current entry bytes. This is intended to be used in hooks
// That will be applied after calling Log().
// The retuned []byte is not a copy and must not be modified directly.
func (e Entry) Bytes() (data []byte) {
	return e.enc.data
}

// Error adds the given error key/value to the log entry
func (e Entry) Error(key string, err error) (entry Entry) {
	if e.enc != nil {
		e.enc.addKey(key)
		e.enc.AppendString(err.Error())
	}
	return e
}

// Bool adds the given bool key/value to the log entry
func (e Entry) Bool(key string, value bool) (entry Entry) {
	if e.enc != nil {
		e.enc.addKey(key)
		e.enc.AppendBool(value)
	}
	return e
}

// Float adds the given float key/value to the log entry
func (e Entry) Float(key string, value float64) (entry Entry) {
	if e.enc != nil {
		e.enc.addKey(key)
		e.enc.AppendFloat(value)
	}
	return e
}

// Int adds the given int key/value to the log entry
func (e Entry) Int(key string, value int64) (entry Entry) {
	if e.enc != nil {
		e.enc.addKey(key)
		e.enc.AppendInt(value)
	}
	return e
}

// Uint adds the given uint key/value to the log entry
func (e Entry) Uint(key string, value uint64) (entry Entry) {
	if e.enc != nil {
		e.enc.addKey(key)
		e.enc.AppendUint(value)
	}
	return e
}

// String adds the given string key/value to the log entry
func (e Entry) String(key string, value string) (entry Entry) {
	if e.enc != nil {
		e.enc.addKey(key)
		e.enc.AppendString(value)
	}
	return e
}

// Object creates a json object
func (e Entry) Object(key string, fn func(Entry)) (entry Entry) {
	if e.enc != nil {
		e.enc.addKey(key)
		e.enc.openObject()
	}

	fn(e)

	if e.enc != nil {
		e.enc.closeObject()
	}
	return e
}

func (e Entry) init(level Level) {

	t := time.Now()
	e.level = level

	e.enc.openObject()
	e.enc.addKey("level")
	e.enc.AppendString(level.String())

	if e.l.config.EnableTime {
		e.enc.addKey(e.l.config.TimeField)

		switch e.l.config.TimeFormat {
		case Unix:
			e.enc.AppendInt(t.Unix())
		case UnixMilli:
			e.enc.AppendInt(t.UnixNano() / 1000000)
		case UnixNano:
			e.enc.AppendInt(t.UnixNano())
		default:
			e.enc.data = append(e.enc.data, '"')
			e.enc.data = t.AppendFormat(e.enc.data, e.l.config.TimeFormat)
			e.enc.data = append(e.enc.data, '"')
		}

	}

	if e.l.config.EnableCaller {
		_, f, l, ok := runtime.Caller(3 + e.l.config.CallerSkip)
		e.enc.addKey("caller")
		if ok {
			idx := strings.LastIndexByte(f, '/')
			idx = strings.LastIndexByte(f[:idx], '/')
			e.enc.AppendString(f[idx+1:] + ":" + strconv.Itoa(l))
		} else if !ok {
			e.enc.AppendString("???")
		}

	}
}
