package log

import (
	logrus "github.com/sirupsen/logrus"
)

func Trace(args ...any) {
	logrus.Trace(args...)
}

func Tracef(format string, args ...any) {
	logrus.Tracef(format, args...)
}

func Debug(args ...any) {
	logrus.Debug(args...)
}

func Debugf(format string, args ...any) {
	logrus.Debugf(format, args...)
}

func Warn(args ...any) {
	logrus.Warn(args...)
}

func Warnf(format string, args ...any) {
	logrus.Warnf(format, args...)
}

func Error(args ...any) {
	logrus.Error(args...)
}

func Errorf(format string, args ...any) {
	logrus.Errorf(format, args...)
}

func Info(args ...any) {
	logrus.Info(args...)
}

func Infof(format string, args ...any) {
	logrus.Infof(format, args...)
}

func Fatal(args ...any) {
	logrus.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	logrus.Fatalf(format, args...)
}

func Panic(args ...any) {
	logrus.Panic(args...)
}

func Panicf(format string, args ...any) {
	logrus.Panicf(format, args...)
}
