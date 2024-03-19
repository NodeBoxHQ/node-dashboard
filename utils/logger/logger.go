package logger

import (
	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.New()
	Logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Trace(args ...interface{}) {
	Logger.Trace(args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}
