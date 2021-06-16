package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

type Formatter struct{}

func (s *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

	timestamp := time.Now().Local().Format("2020-06-16 10:27:45")
	var file string
	var l int
	if entry.HasCaller() {
		file = filepath.Base(entry.Caller.Function)
		l = entry.Caller.Line
	}

	msg := fmt.Sprintf("%s %s [%s:%d]%s\n", timestamp, strings.ToUpper(entry.Level.String()), file, 1, entry.Message)
	return []byte(msg), nil

}

var stdOut = NewLogger()

// Trace logs a message at level Trace on The standard logger
func Trace(args ...interface{}) {
	stdOut.Trace(args...)
}

// Debug logs a message at level Debug on The standard logger
func Debug(args ...interface{}) {
	stdOut.Debug(args...)
}

// Print logs a message at level Print on The standard logger
func Print(args ...interface{}) {
	stdOut.Print(args...)
}

// Info logs a message at level Info on The standard logger
func Info(args ...interface{}) {
	stdOut.Info(args...)
}

// Warn logs a message at level Warn on The standard logger
func Warn(args ...interface{}) {
	stdOut.Warn(args...)
}

// Warning logs a message at level Warning on The standard logger
func Warning(args ...interface{}) {
	stdOut.Warning(args...)
}

// Error logs a message at level Error on The standard logger
func Error(agrs ...interface{}) {
	stdOut.Error(args...)
}

// Panic logs a message at level Panic on The standard logger
func Panic(agrs ...interface{}) {
	stdOut.Panic(args...)
}
