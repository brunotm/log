// Package log provides a simple, leveled, fast, zero allocation, json structured logging package for Go
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
	"io"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// ISO8601 time format
	ISO8601 = "2006-01-02T15:04:05.999Z07:00"
	// Unix time in seconds
	Unix = "unix"
	// UnixMilli time in milliseconds
	UnixMilli = "unix_milli"
	// UnixNano time in nanoseconds
	UnixNano = "unix_nano"

	entrySize = 512
)

var (
	encoderPool *sync.Pool

	// DefaultConfig for logger
	DefaultConfig = Config{
		Format:         FormatJSON,
		Level:          INFO,
		EnableCaller:   true,
		CallerSkip:     0,
		EnableTime:     true,
		TimeField:      "time",
		TimeFormat:     ISO8601,
		MessageField:   "message",
		LevelField:     "level",
		EnableSampling: true,
		SamplingTick:   time.Second,
		SamplingStart:  100,
		SamplingFactor: 100,
	}
)

func init() {
	encoderPool = &sync.Pool{
		New: newEncoder,
	}

	for x := 0; x < 32; x++ {
		encoderPool.Put(newEncoder())
	}
}

func newEncoder() interface{} {
	return &encoder{data: make([]byte, 0, entrySize)}
}

// Config type for logger
type Config struct {
	Format         Format        // Log format
	Level          Level         // Log level
	EnableCaller   bool          // Enable caller info
	CallerSkip     int           // Skip level of callers, useful if wrapping the logger
	EnableTime     bool          // Enable log timestamps
	TimeField      string        // Field name for the log timestamp
	TimeFormat     string        // Time Format for log timestamp
	MessageField   string        // Field name for the log message
	LevelField     string        // Field name for the log level
	EnableSampling bool          // Enable log sampling to reduce CPU and I/O load
	SamplingTick   time.Duration // Resolution at which entries will be sampled
	SamplingStart  int           // Start sampling after this number of similar entries within SamplingTick
	SamplingFactor int           // Reduction factor when sampling
}

// Logger type
type Logger struct {
	config  *Config
	writer  io.Writer
	hooks   []func(Entry)
	with    []func(Entry)
	sampler *sampler
}

// New creates a new logger with the give config and writer.
// A nill writer will be set to ioutil.Discard.
func New(writer io.Writer, config Config) (logger *Logger) {

	if writer == nil {
		writer = ioutil.Discard
	}

	logger = &Logger{}

	if config.EnableSampling {
		logger.sampler = newSampler(
			config.SamplingTick,
			config.SamplingStart,
			config.SamplingFactor)
	}

	logger.writer = writer
	logger.config = &config

	return logger
}

// SetLevel atomically sets the new log level
func (l *Logger) SetLevel(lv Level) {
	atomic.StoreUint32((*uint32)(&l.config.Level), uint32(lv))
}

// SetFormat atomically sets the new log format
func (l *Logger) SetFormat(f Format) {
	atomic.StoreUint32((*uint32)(&l.config.Format), uint32(f))
}

// With creates a new logger with functions to apply context to the log entries.
// With functions are cumulative and applied before all other log data.
func (l *Logger) With(f ...func(Entry)) (logger *Logger) {
	return &Logger{
		config:  l.config,
		writer:  l.writer,
		hooks:   l.hooks,
		sampler: l.sampler,
		with:    append(l.with, f...),
	}
}

// Hooks creates a new logger with functions to apply after the entry is written.
// Hooks are cumulative and useful for shipping log data to other systems.
func (l *Logger) Hooks(f ...func(Entry)) (logger *Logger) {
	return &Logger{
		config:  l.config,
		writer:  l.writer,
		hooks:   append(l.hooks, f...),
		with:    l.with,
		sampler: l.sampler,
	}
}

// entry creates a new log entry with the specified level to be manipulated directly
func (l *Logger) entry(level Level, message string) (entry Entry) {
	entry.level = level

	// Only initialize Entry if on or above the logger Level
	if entry.level >= Level(atomic.LoadUint32((*uint32)(&l.config.Level))) {

		if l.config.EnableSampling && !l.sampler.check(level, message) {
			return entry
		}

		entry.o.enc = encoderPool.Get().(*encoder)
		entry.o.enc.format = Format(atomic.LoadUint32((*uint32)(&l.config.Format)))

		entry.l = l
		entry.init(level)

		for i := 0; i < len(l.with); i++ {
			l.with[i](entry)
		}

		entry.String(l.config.MessageField, message)
	}

	return entry
}

// Debug creates a new log entry with the given message.
func (l *Logger) Debug(message string) (entry Entry) {
	entry = l.entry(DEBUG, message)
	return entry
}

// Info creates a new log entry with the given message.
func (l *Logger) Info(message string) (entry Entry) {
	entry = l.entry(INFO, message)
	return entry
}

// Warn creates a new log entry with the given message.
func (l *Logger) Warn(message string) (entry Entry) {
	entry = l.entry(WARN, message)
	return entry
}

// Error creates a new log entry with the given message.
func (l *Logger) Error(message string) (entry Entry) {
	entry = l.entry(ERROR, message)
	return entry
}

// Fatal creates a new log entry with the given message.
// After write, Fatal calls os.Exit(1) terminating the running program
func (l *Logger) Fatal(message string) (entry Entry) {
	entry = l.entry(FATAL, message)
	return entry
}

func (l *Logger) write(entry Entry) {
	if entry.o.enc == nil {
		return
	}

	defer l.discard(entry)
	l.writer.Write(append(entry.o.enc.data, '\n'))
}

func (l *Logger) discard(entry Entry) {
	for i := 0; i < len(l.hooks); i++ {
		l.hooks[i](entry)
	}

	if entry.level == FATAL {
		os.Exit(1)
	}

	entry.o.enc.reset()
	encoderPool.Put(entry.o.enc)
}
