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

func TestLogObject(t *testing.T) {
	config := DefaultConfig
	config.EnableTime = false
	config.EnableCaller = false
	config.Level = DEBUG
	l := New(os.Stdout, config)

	l.Hooks(func(e Entry) {
		w := []byte(`{"level":"error","message":"error message","string value":"text","int value":8,"object":{"user":"userA","id":72386784}}`)
		if bytes.Equal(w, e.Bytes()) {
			t.Fatal("error logging object")
		}
	})

	l.Error("error message").
		String("string value", "text").
		Int("int value", 8).
		Object("object", func(o Object) {
			o.String("user", "userA").Int("id", 72386784)
		}).Write()
}

func TestLogArray(t *testing.T) {
	config := DefaultConfig
	config.EnableTime = false
	config.EnableCaller = false
	config.Level = DEBUG
	l := New(os.Stdout, config)

	l.Hooks(func(e Entry) {
		w := []byte(`{"level":"error","message":"error message","string value":"text","int value":8,"array":["userA",72386784,null]}`)
		if !bytes.Equal(w, e.Bytes()) {
			t.Fatal("error logging object")
		}
	})

	l.Error("error message").
		String("string value", "text").
		Int("int value", 8).
		Array("array", func(a Array) {
			a.AppendString("userA").AppendInt(72386784).AppendNull()
		}).Write()
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
			Int("int value", 8).Float("float", 722727272.0099).
			Int("int", 8).Float("float value", 722727272.0099).
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
			Int("int value", 8).Float("float", 722727272.0099).
			Int("int", 8).Float("float value", 722727272.0099).
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
			Int("int value", 8).Float("float", 722727272.0099).
			Int("int", 8).Float("float value", 722727272.0099).
			Write()
	}
}

func Example() {
	config := DefaultConfig
	config.Level = DEBUG

	// New logger with added context
	l := New(os.Stderr, config).
		With(func(e Entry) {
			e.String("application", "app1")
		})

	// Simple logging
	l.Info("info message").String("key", "value").Write()
	// {"level":"info","time":"2019-01-30T20:42:56.445Z","caller":"_local/main.go:21",
	// "application":"app1","message":"info message","key":"value"}

	l.Warn("warn message").Bool("flag", false).Write()
	// {"level":"warn","time":"2019-01-30T20:42:56.446Z","caller":"_local/main.go:24",
	// "application":"app1","message":"warn message","flag":false}

	l.Error("caught an error").String("error", "request error").Write()
	// {"level":"error","time":"2019-01-30T20:42:56.446Z","caller":"_local/main.go:27",
	// "application":"app1","message":"caught an error","error":"request error"}

	// Create nested objects in log entry
	l.Debug("debug message").Object("request_data", func(o Object) {
		o.String("request_id", "4BA0D8B1-4ABA-4D70-A55C-3358667C058B").
			String("user_id", "3B1BA12B-68DF-4DB7-809B-1AC5D8AF663A").
			Float("value", 3.1415926535)
	}).Write()
	// {"level":"debug","time":"2019-01-30T22:44:45.193Z","caller":"_local/main.go:31",
	// "application":"app1","message":"debug message","request_data":
	// {"request_id":"4BA0D8B1-4ABA-4D70-A55C-3358667C058B",
	// "user_id":"3B1BA12B-68DF-4DB7-809B-1AC5D8AF663A","value":3.1415926535}}

	// Create array objects in log entry
	l.Debug("debug message").Array("request_points", func(a Array) {
		a.AppendFloat(3.1415926535).
			AppendFloat(2.7182818284).
			AppendFloat(1.41421).
			AppendFloat(1.6180339887498948482)
	}).Write()
	// {"level":"debug","time":"2019-02-04T08:42:15.216Z","caller":"_local/main.go:44",
	// "application":"app1","message":"debug message",
	// "request_points":[3.1415926535,2.7182818284,1.41421,1.618033988749895]}

}

type writerCounter struct {
	count int
}

func (w *writerCounter) Write(p []byte) (n int, err error) {
	w.count++
	return len(p), nil
}
