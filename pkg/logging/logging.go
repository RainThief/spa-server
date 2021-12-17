package logging

import (
	"errors"
	"fmt"
	"io"
)

const (
	// FatalLevel only
	FatalLevel = iota
	// ErrorLevel and below
	ErrorLevel
	// WarningLevel and below
	WarningLevel
	// DebugLevel and below
	DebugLevel
	// InfoLevel is all events
	InfoLevel
)

func setLogLevel(level string) (logLevel int) {
	switch level {
	case "FATAL":
		logLevel = FatalLevel
	case "ERROR":
		logLevel = ErrorLevel
	case "WARNING":
		logLevel = WarningLevel
	case "DEBUG":
		logLevel = DebugLevel
	case "INFO":
		logLevel = InfoLevel
	default:
		logLevel = InfoLevel
	}

	return
}

// Logger provides logging functionality
type Logger struct {
	stdOut   io.Writer
	stdErr   io.Writer
	logLevel int
}

// NewLogger returns logger instance
func NewLogger(stdOut, stdErr io.Writer, level string) *Logger {
	return &Logger{stdOut, stdErr, setLogLevel(level)}
}

// Fatal logs fatal error event
func (l *Logger) Fatal(msg string, v ...interface{}) {
	_, _ = l.stdErr.Write([]byte("FATAL: " + fmt.Sprintf(msg+"\n", v...)))
}

// Error logs error event
func (l *Logger) Error(msg string, v ...interface{}) {
	if l.logLevel >= ErrorLevel {
		_, _ = l.stdErr.Write([]byte("ERROR: " + fmt.Sprintf(msg+"\n", v...)))
	}
}

// Warning log event
func (l *Logger) Warning(msg string, v ...interface{}) {
	if l.logLevel >= WarningLevel {
		_, _ = l.stdErr.Write([]byte("WARNING: " + fmt.Sprintf(msg+"\n", v...)))
	}
}

// Debug logs debug event
func (l *Logger) Debug(msg string, v ...interface{}) {
	if l.logLevel >= DebugLevel {
		_, _ = l.stdOut.Write([]byte("DEBUG: " + fmt.Sprintf(msg+"\n", v...)))
	}
}

// Info logs info event
func (l *Logger) Info(msg string, v ...interface{}) {
	if l.logLevel >= InfoLevel {
		_, _ = l.stdOut.Write([]byte("INFO: " + fmt.Sprintf(msg+"\n", v...)))
	}
}

// LogAndRaiseError logs error event and raises error
func (l *Logger) LogAndRaiseError(msg string, v ...interface{}) error {
	errorMsg := fmt.Sprintf(msg, v...)
	if l.logLevel >= ErrorLevel {
		l.Error(errorMsg)
	}

	return errors.New(errorMsg)
}

// SetSTDOutput provides realtime changing of output writer for stdout
func (l *Logger) SetSTDOutput(w io.Writer) {
	l.stdOut = w
}

// SetSTDError provides realtime changing of output writer for stderr
func (l *Logger) SetSTDError(w io.Writer) {
	l.stdErr = w
}
