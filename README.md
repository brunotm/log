# log

[![Build Status](https://travis-ci.org/brunotm/log.svg?branch=master)](https://travis-ci.org/brunotm/log)
[![Go Report Card](https://goreportcard.com/badge/github.com/brunotm/log)](https://goreportcard.com/report/github.com/brunotm/log)
[![GoDoc](https://godoc.org/github.com/brunotm/log?status.svg)](https://godoc.org/github.com/brunotm/log)

A simple, leveled, fast, zero allocation, structured logging package for Go.
Designed to make logging on the hot path dirt cheap, dependency free and my life easier.

A simpler textual key/value format is also available by setting `Config.Text to true`

By default log entries with the same level and message within each second, will be sampled to cap CPU and I/O load under high logging activity.
This behavior can be disabled by setting `Config.EnableSampling to false`

## Usage

```go
package main

import (
    "os"

    "github.com/brunotm/log"
)

func main() {
    config := DefaultConfig
    config.Level = DEBUG
    config.Text = true  // Enable text logging

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
```

## Performance on a 2,3 GHz Intel Core i5, 2017 13-inch Macbook Pro

Message: `{"level":"info","time":"2019-01-30T20:54:07.029Z","message":"informational message","string value":"text","int value":8,"float":722727272.0099,"int":8,"float value":722727272.0099}`

```bash
pkg: github.com/brunotm/log
BenchmarkLogNoSampling-8         1672147               704 ns/op               0 B/op          0 allocs/op
BenchmarkLogWithSampling-8      10266382               111 ns/op               0 B/op          0 allocs/op
BenchmarkLogNoLevel-8           340966922                3.46 ns/op            0 B/op          0 allocs/op
```

## Contact

Bruno Moura [brunotm@gmail.com](mailto:brunotm@gmail.com)

## License

log source code is available under the Apache Version 2.0 [License](/LICENSE)
