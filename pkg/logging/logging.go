package logging

import (
	"errors"
	"fmt"
	"io"
)

// Logger provides logging functionality
type Logger struct {
	stdOut io.Writer
	stdErr io.Writer
}

// @todo use log package and set loglevel

// Error logs error event
func (l *Logger) Error(msg string, v ...interface{}) {
	_, _ = l.stdErr.Write([]byte("ERROR: " + fmt.Sprintf(msg+"\n", v...)))
}

// Fatal logs fatal error event
func (l *Logger) Fatal(msg string, v ...interface{}) {
	_, _ = l.stdErr.Write([]byte("FATAL: " + fmt.Sprintf(msg+"\n", v...)))
}

// Debug logs debug event
func (l *Logger) Debug(msg string, v ...interface{}) {
	_, _ = l.stdOut.Write([]byte("DEBUG: " + fmt.Sprintf(msg+"\n", v...)))
}

// Info logs info event
func (l *Logger) Info(msg string, v ...interface{}) {
	_, _ = l.stdOut.Write([]byte("INFO: " + fmt.Sprintf(msg+"\n", v...)))
}

// LogAndRaiseError logs error event and raises error
func (l *Logger) LogAndRaiseError(msg string, v ...interface{}) error {
	errorMsg := fmt.Sprintf(msg, v...)
	l.Error(errorMsg)
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
