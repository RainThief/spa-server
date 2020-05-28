package logging

import (
	"fmt"
	"os"
)

// @todo use log package

func Error(msg string, v ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(msg+"\n", v...))
}

func Fatal(msg string, v ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(msg+"\n", v...))
}

func Debug(msg string, v ...interface{}) {
	os.Stdout.WriteString(fmt.Sprintf(msg+"\n", v...))
}

func Info(msg string, v ...interface{}) {
	os.Stdout.WriteString(fmt.Sprintf(msg+"\n", v...))
}
