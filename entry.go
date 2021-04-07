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

// Entry is a structured log entry. A entry is not safe for concurrent use.
// A entry must be logged by calling Log(), and cannot be reused after.
type Entry struct {
	o     Object
	l     *Logger
	level Level
}

// Write logs the current entry. An entry must not be used after calling Write().
func (e Entry) Write() {
	if e.o.enc != nil {
		if e.o.enc.format == FormatJSON {
			e.o.enc.closeObject()
		}
		e.l.write(e)
	}
}

// Level returns the log level of current entry.
func (e Entry) Level() (level Level) {
	return e.level
}

// Bytes return the current entry bytes. This is intended to be used in hooks
// That will be applied after calling Log().
// The returned []byte is not a copy and must not be modified directly.
func (e Entry) Bytes() (data []byte) {
	return e.o.enc.data
}

// Bool adds the given bool key/value
func (e Entry) Bool(key string, value bool) (entry Entry) {
	if e.o.enc != nil {
		e.o.Bool(key, value)
	}
	return e
}

// Float32 adds the given float32 key/value
func (e Entry) Float32(key string, value float32) (entry Entry) {
	e.Float64(key, float64(value))
	return e
}

// Float64 adds the given float64 key/value
func (e Entry) Float64(key string, value float64) (entry Entry) {
	if e.o.enc != nil {
		e.o.Float64(key, value)
	}
	return e
}

// Int8 adds the given int8 key/value
func (e Entry) Int8(key string, value int8) (entry Entry) {
	e.Int64(key, int64(value))
	return e
}

// Int16 adds the given int16 key/value
func (e Entry) Int16(key string, value int16) (entry Entry) {
	e.Int64(key, int64(value))
	return e
}

// Int32 adds the given int32 key/value
func (e Entry) Int32(key string, value int32) (entry Entry) {
	e.Int64(key, int64(value))
	return e
}

// Int adds the given int key/value
func (e Entry) Int(key string, value int) (entry Entry) {
	e.Int64(key, int64(value))
	return e
}

// Int64 adds the given int64 key/value
func (e Entry) Int64(key string, value int64) (entry Entry) {
	if e.o.enc != nil {
		e.o.Int64(key, value)
	}
	return e
}

// Uint8 adds the given uint8 key/value
func (e Entry) Uint8(key string, value uint8) (entry Entry) {
	e.Uint64(key, uint64(value))
	return e
}

// Uint16 adds the given uint16 key/value
func (e Entry) Uint16(key string, value uint16) (entry Entry) {
	e.Uint64(key, uint64(value))
	return e
}

// Uint32 adds the given uint32 key/value
func (e Entry) Uint32(key string, value uint32) (entry Entry) {
	e.Uint64(key, uint64(value))
	return e
}

// Uint adds the given uint16 key/value
func (e Entry) Uint(key string, value uint) (entry Entry) {
	e.Uint64(key, uint64(value))
	return e
}

// Uint64 adds the given uint key/value
func (e Entry) Uint64(key string, value uint64) (entry Entry) {
	if e.o.enc != nil {
		e.o.Uint64(key, value)
	}
	return e
}

// String adds the given string key/value
func (e Entry) String(key string, value string) (entry Entry) {
	if e.o.enc != nil {
		e.o.String(key, value)
	}
	return e
}

// Null adds a null value for the given key
func (e Entry) Null(key string) (entry Entry) {
	if e.o.enc != nil {
		e.o.Null(key)
	}
	return e
}

// Error adds the given error key/value
func (e Entry) Error(key string, value error) (entry Entry) {
	if e.o.enc != nil {
		e.o.Error(key, value)
	}
	return e
}

func (e Entry) init(level Level) {

	t := time.Now()
	e.level = level

	if e.l.config.EnableTime {
		e.o.enc.addKey(e.l.config.TimeField)

		switch e.l.config.TimeFormat {
		case Unix:
			e.o.enc.AppendInt64(t.Unix())
		case UnixMilli:
			e.o.enc.AppendInt64(t.UnixNano() / int64(time.Millisecond))
		case UnixNano:
			e.o.enc.AppendInt64(t.UnixNano())
		default:
			e.o.enc.data = append(e.o.enc.data, '"')
			e.o.enc.data = t.AppendFormat(e.o.enc.data, e.l.config.TimeFormat)
			e.o.enc.data = append(e.o.enc.data, '"')
		}

	}

	e.String(e.l.config.LevelField, level.String())

	if e.l.config.EnableCaller {
		_, f, l, ok := runtime.Caller(3 + e.l.config.CallerSkip)

		if ok {
			idx := strings.LastIndexByte(f, '/')
			idx = strings.LastIndexByte(f[:idx], '/')
			e.String("caller", f[idx+1:]+":"+strconv.Itoa(l))
		} else {
			e.String("caller", "???")
		}
	}
}
