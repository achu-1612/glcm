package log

import (
	"io"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

var std *log.Logger
var mu *sync.Mutex

func init() {
	std = log.New()

	std.SetOutput(os.Stdout)
	std.SetLevel(log.InfoLevel)

	std.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
		ForceColors:   true,
	})
}

func SetOutput(o io.Writer) {
	mu.Lock()
	defer mu.Lock()

	std.SetOutput(o)
}

func Info(args ...interface{}) {
	std.Info(args...)
}

func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

func Debug(args ...interface{}) {
	std.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	std.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

func Error(args ...interface{}) {
	std.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}
