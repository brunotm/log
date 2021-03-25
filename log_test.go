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
	"errors"
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
			w := []byte(`{"level":"debug","message":"debug message","string value":"text","int value":8,"null":null,"error":"new error",}`)
			if !bytes.Equal(w, e.Bytes()) {
				t.Fatal("error logging warn")
			}
		case INFO:
			w := []byte(`{"level":"info","message":"info message","string value":"text","int value":8,"null value":null,"error":"new error",}`)
			if !bytes.Equal(w, e.Bytes()) {
				t.Fatal("error logging info")
			}
		case WARN:
			w := []byte(`{"level":"warn","message":"warn message","string value":"text","int value":8,"null value":null,"error":"new error",}`)
			if !bytes.Equal(w, e.Bytes()) {
				t.Fatal("error logging warn")
			}
		case ERROR:
			w := []byte(`{"level":"error","message":"error message","string value":"text","int value":8,"null value":null,"error":"new error",}`)
			if !bytes.Equal(w, e.Bytes()) {
				t.Fatal("error logging error")
			}
		}
	})

	l.Debug("debug message").
		String("string value", "text").
		Int("int value", 8).Null("null value").
		Error("error", errors.New("new error")).Write()

	l.Info("info message").
		String("string value", "text").
		Int("int value", 8).Null("null value").
		Error("error", errors.New("new error")).Write()

	l.Warn("warn message").
		String("string value", "text").
		Int("int value", 8).Null("null value").
		Error("error", errors.New("new error")).Write()

	l.Error("error message").
		String("string value", "text").
		Int("int value", 8).Null("null value").
		Error("error", errors.New("new error")).Write()

}

func TestLogText(t *testing.T) {
	config := DefaultConfig
	config.EnableTime = false
	config.EnableCaller = false
	config.Level = DEBUG
	config.Text = true
	l := New(os.Stdout, config)

	l.Hooks(func(e Entry) {
		w := []byte(`level="debug" message="debug message" string="text" int=8 null=null error="new error"`)
		if !bytes.Equal(w, e.Bytes()) {
			t.Fatal("error logging warn")
		}
	})

	l.Debug("debug message").
		String("string", "text").
		Int("int", 8).Null("null").
		Error("error", errors.New("new error")).Write()
}

func TestLogSampler(t *testing.T) {
	logCount := 10000
	w := &writerCounter{}
	config := DefaultConfig

	l := New(w, config)

	for x := 0; x < logCount; x++ {
		l.Error("error message").
			String("string value", "text").
			Int("int value", 8).
			Write()
	}

	if logCount <= w.count {
		t.Fatalf("number of interactions %d number of writes %d", logCount, w.count)
	}

}

func BenchmarkLogNoSampling(b *testing.B) {
	config := DefaultConfig
	config.Level = DEBUG
	config.EnableCaller = false
	config.EnableSampling = false

	l := New(ioutil.Discard, config)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("informational message").
			String("string value", "text").
			Int("int value", 8).Float64("float", 722727272.0099).
			Int("int", 8).Float64("float value", 722727272.0099).
			Write()
	}
}

func BenchmarkLogWithSampling(b *testing.B) {
	config := DefaultConfig
	config.Level = DEBUG
	config.EnableCaller = false
	config.EnableSampling = true

	l := New(ioutil.Discard, config)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("informational message").
			String("string value", "text").
			Int("int value", 8).Float64("float", 722727272.0099).
			Int("int", 8).Float64("float value", 722727272.0099).
			Write()
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
			Int("int value", 8).Float64("float", 722727272.0099).
			Int("int", 8).Float64("float value", 722727272.0099).
			Write()
	}
}

func Example() {
	config := DefaultConfig
	config.Level = DEBUG
	config.Text = true // Enable text logging

	// New logger with added context
	l := New(os.Stdout, config).
		With(func(e Entry) {
			e.String("app", "app1")
		})

	// Simple logging
	l.Info("info message").String("key", "value").Write()
	// Text format:
	// time="2021-03-25T13:32:50.391Z" level="info" caller="_local/main.go:26" app="app1" message="info message" key="value"
	// JSON format:
	// {"time":"2021-03-25T13:33:20.547Z", "level":"info", "caller":"_local/main.go:26", "app":"app1", "message":"info message", "key":"value"}

	l.Warn("warn message").Bool("flag", false).Write()
	// Text format:
	// time="2021-03-25T13:32:50.391Z" level="warn" caller="_local/main.go:29" app="app1" message="warn message" flag=false
	// JSON format:
	// {"time":"2021-03-25T13:33:20.547Z", "level":"warn", "caller":"_local/main.go:29", "app":"app1", "message":"warn message", "flag":false}

	l.Error("caught an error").String("error", "request error").Write()
	// Text format:
	// time="2021-03-25T13:32:50.391Z" level="error" caller="_local/main.go:32" app="app1" message="caught an error" error="request error"
	// JSON format:
	// {"time":"2021-03-25T13:33:20.547Z", "level":"error", "caller":"_local/main.go:32", "app":"app1", "message":"caught an error", "error":"request error"}

	l.Fatal("caught an unrecoverable error").Error("error", errors.New("some error")).Write()
	// Text format:
	// time="2021-03-25T13:32:50.391Z" level="fatal" caller="_local/main.go:35" app="app1" message="caught an unrecoverable error" error="some error"
	// JSON format:
	// {"time":"2021-03-25T13:33:20.547Z", "level":"fatal", "caller":"_local/main.go:35", "app":"app1", "message":"caught an unrecoverable error", "error":"some error"}
}

type writerCounter struct {
	count int
}

func (w *writerCounter) Write(p []byte) (n int, err error) {
	w.count++
	return len(p), nil
}
