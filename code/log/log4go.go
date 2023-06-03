package log

import (
	"math"
	"path"

	log "github.com/meission/log4go"
)

var (
	logger log.Logger
)

type Config struct {
	Dir string
}

func Init(c *Config) {
	if c == nil || c.Dir == "" {
		c = &Config{
			Dir: "./",
		}
	}
	logger = log.Logger{}
	log.LogBufferLength = 10240
	// new info writer
	iw := log.NewFileLogWriter(path.Join(c.Dir, "info.log"), false)
	iw.SetRotateDaily(true)
	iw.SetRotateSize(math.MaxInt32)
	iw.SetRotate(true)
	iw.SetFormat("[%D %T] [%L] [%S] %M")
	logger["info"] = &log.Filter{
		Level:     log.INFO,
		LogWriter: iw,
	}
	// new warning writer
	ww := log.NewFileLogWriter(path.Join(c.Dir, "warning.log"), false)
	ww.SetRotateDaily(true)
	ww.SetRotateSize(math.MaxInt32)
	ww.SetRotate(true)
	ww.SetFormat("[%D %T] [%L] [%S] %M")
	logger["warning"] = &log.Filter{
		Level:     log.WARNING,
		LogWriter: ww,
	}
	// new error writer
	ew := log.NewFileLogWriter(path.Join(c.Dir, "error.log"), false)
	ew.SetRotateDaily(true)
	ew.SetRotateSize(math.MaxInt32)
	ew.SetRotate(true)
	ew.SetFormat("[%D %T] [%L] [%S] %M")
	logger["error"] = &log.Filter{
		Level:     log.ERROR,
		LogWriter: ew,
	}
}

// Close close resource.
func Close() {
	if logger != nil {
		logger.Close()
	}
}

// Info write info log .
func Info(format string, args ...interface{}) {
	if logger != nil {
		logger.Info(format, args...)
	}
}

// Warn write warn log .
func Warn(format string, args ...interface{}) {
	if logger != nil {
		logger.Warn(format, args...)
	}
}

// Error write error log .
func Error(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(format, args...)
	}
}
