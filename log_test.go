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
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestLogEntry(t *testing.T) {
	config := DefaultConfig
	config.EnableTime = false
	config.EnableCaller = false
	config.Level = DEBUG
	l := New(os.Stdout, config)

	l.Hooks(func(e Entry) {
		switch e.Level() {
		case DEBUG:
			if 0 != bytes.Compare(
				[]byte(`{"level":"debug","message":"debug message","string value":"text","int value":8}`), e.Bytes()) {
				t.Fatal("error logging warn")
			}
		case INFO:
			if 0 != bytes.Compare(
				[]byte(`{"level":"info","message":"info message","string value":"text","int value":8}`), e.Bytes()) {
				t.Fatal("error logging info")
			}
		case WARN:
			if 0 != bytes.Compare(
				[]byte(`{"level":"warn","message":"warn message","string value":"text","int value":8}`), e.Bytes()) {
				t.Fatal("error logging warn")
			}
		case ERROR:
			if 0 != bytes.Compare(
				[]byte(`{"level":"error","message":"error message","string value":"text","int value":8}`), e.Bytes()) {
				t.Fatal("error logging error")
			}
		}
	})

	l.Debug("debug message").
		String("string value", "text").
		Int("int value", 8).Log()

	l.Info("info message").
		String("string value", "text").
		Int("int value", 8).Log()

	l.Warn("warn message").
		String("string value", "text").
		Int("int value", 8).Log()

	l.Error("error message").
		String("string value", "text").
		Int("int value", 8).Log()

}

func TestLogObject(t *testing.T) {
	config := DefaultConfig
	config.EnableTime = false
	config.EnableCaller = false
	config.Level = DEBUG
	l := New(os.Stdout, config)

	l.Hooks(func(e Entry) {
		if 0 != bytes.Compare([]byte(`{"level":"error","message":"error message","string value":"text","int value":8,"object":{"user":"userA","id":72386784}}`), e.Bytes()) {
			t.Fatal("error logging object")
		}
	})

	l.Error("error message").
		String("string value", "text").
		Int("int value", 8).
		Object("object", func(e Entry) {
			e.String("user", "userA").Int("id", 72386784)
		}).Log()
}

func TestLogArray(t *testing.T) {
	config := DefaultConfig
	config.EnableTime = false
	config.EnableCaller = false
	config.Level = DEBUG
	l := New(os.Stdout, config)

	l.Hooks(func(e Entry) {
		if 0 != bytes.Compare([]byte(`{"level":"error","message":"error message","string value":"text","int value":8,"array":["userA",72386784]}`), e.Bytes()) {
			t.Fatal("error logging object")
		}
	})

	l.Error("error message").
		String("string value", "text").
		Int("int value", 8).
		Array("array", func(a Array) {
			a.AppendString("userA").AppendInt(72386784)
		}).Log()
}

func BenchmarkLog(b *testing.B) {
	config := DefaultConfig
	config.Level = DEBUG
	config.EnableCaller = false
	l := New(ioutil.Discard, config)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("informational message").
			String("string value", "text").
			Int("int value", 8).Float("float", 722727272.0099).
			Int("int", 8).Float("float value", 722727272.0099).
			Log()
	}
}

func BenchmarkLogNoLevel(b *testing.B) {
	config := DefaultConfig
	config.Level = ERROR
	config.EnableCaller = false
	l := New(os.Stdout, config)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("informational message").
			String("string value", "text").
			Int("int value", 8).Float("float", 722727272.0099).
			Int("int", 8).Float("float value", 722727272.0099).
			Log()
	}
}
