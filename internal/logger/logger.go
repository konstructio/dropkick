package logger

import (
	"io"
	"log"
)

// Logger is a struct that wraps the standard library's log.Logger
type Logger struct {
	log *log.Logger
}

// None is a default Logger instance with no underlying log.Logger
var None = &Logger{}

// New creates a new Logger instance. It takes an io.Writer as the destination
// where it will write log entries. The log entries will be prefixed with the date and time.
func New(dest io.Writer) *Logger {
	return &Logger{log: log.New(dest, "", log.Ldate|log.Ltime)}
}

// Printf formats according to a format specifier and writes to the logger's io.Writer.
// If the logger's underlying log.Logger is nil, it does nothing.
func (l *Logger) Printf(format string, v ...interface{}) {
	if l.log == nil {
		return
	}

	l.log.Printf(format, v...)
}
