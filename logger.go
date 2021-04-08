package log

import (
	"os"
)

var logger *Logger

func init() {
	config := DefaultConfig
	config.CallerSkip++
	logger = New(os.Stderr, config)
}

// SetFormat sets the logging format for the default package logger
func SetFormat(f Format) {
	logger.SetFormat(f)
}

// SetLevel sets the logging level for the default package logger
func SetLevel(l Level) {
	logger.SetLevel(l)
}

// Debug creates a new log entry with the given message with the default package logger.
func Debug(message string) (entry Entry) {
	return logger.Debug(message)
}

// Info creates a new log entry with the given message with the default package logger.
func Info(message string) (entry Entry) {
	return logger.Info(message)
}

// Warn creates a new log entry with the given message with the default package logger.
func Warn(message string) (entry Entry) {
	return logger.Warn(message)
}

// Error creates a new log entry with the given message with the default package logger.
func Error(message string) (entry Entry) {
	return logger.Error(message)
}

// Fatal creates a new log entry with the given message with the default package logger.
// After write, Fatal calls os.Exit(1) terminating the running program
func Fatal(message string) (entry Entry) {
	return logger.Fatal(message)
}
