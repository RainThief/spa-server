package logging

import (
	"fmt"
	"os"
)

// @todo use log package

func Error(msg string, v ...interface{}) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf(msg+"\n", v...))
}

func Fatal(msg string, v ...interface{}) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf(msg+"\n", v...))
}

func Debug(msg string, v ...interface{}) {
	_, _ = os.Stdout.WriteString(fmt.Sprintf(msg+"\n", v...))
}

func Info(msg string, v ...interface{}) {
	_, _ = os.Stdout.WriteString(fmt.Sprintf(msg+"\n", v...))
}
