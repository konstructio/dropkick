package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
)

var (
	colorCyan   = color.New(color.FgCyan).SprintFunc()
	colorYellow = color.New(color.FgYellow).SprintFunc()
	colorRed    = color.New(color.FgRed).SprintFunc()
)

// Logger is a struct that wraps the standard library's log.Logger
type Logger struct {
	log io.Writer
}

// None is a default Logger instance with no underlying log.Logger
var None = &Logger{}

// New creates a new Logger instance. It takes an io.Writer as the destination
// where it will write log entries. The log entries will be prefixed with the date and time.
func New(dest io.Writer) *Logger {
	return &Logger{log: dest}
}

// Infof formats according to a format specifier and writes to the logger's io.Writer.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.write(colorCyan(fmt.Sprintf(format, v...)))
}

// Warnf formats according to a format specifier and writes to the logger's io.Writer.
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.write(colorYellow(fmt.Sprintf(format, v...)))
}

// Errorf formats according to a format specifier and writes to the logger's io.Writer.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.write(colorRed(fmt.Sprintf(format, v...)))
}

// printf is an internal method that writes to the logger's io.Writer or does
// nothing if the logger's io.Writer is nil.
func (l *Logger) write(contents string) {
	if l.log == nil {
		return
	}

	fmt.Fprint(l.log, time.Now().Format("2006/01/02 15:04:05")+" "+contents+"\n")
}
