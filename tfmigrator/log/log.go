package log

import (
	"errors"
	"log"
)

// Logger is a logger.
type Logger interface {
	Info(string)
	Debug(string)
	SetLogLevel(string) error
}

// SimpleLogger is a logger using the standard "log" package.
type SimpleLogger struct {
	logLevel string
}

// Info outputs info log.
func (logger *SimpleLogger) Info(msg string) {
	log.Println("[INFO] " + msg)
}

// Debug outputs debug log.
func (logger *SimpleLogger) Debug(msg string) {
	if logger.logLevel == "debug" {
		log.Println("[DEBUG] " + msg)
	}
}

var logLevels = map[string]struct{}{ //nolint:gochecknoglobals
	"info":  {},
	"debug": {},
}

// SetLogLevel sets the log level
func (logger *SimpleLogger) SetLogLevel(level string) error {
	if _, ok := logLevels[level]; !ok {
		return errors.New("invalid log level: " + level)
	}
	logger.logLevel = level
	return nil
}
