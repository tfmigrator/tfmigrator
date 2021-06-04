package tfmigrator

import "log"

// Logger is a logger.
type Logger interface {
	Info(string)
	Debug(string)
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
