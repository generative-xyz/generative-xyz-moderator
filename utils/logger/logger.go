package logger

import (
	"github.com/sirupsen/logrus"
)

type Ilogger interface {
	Info(fields ...interface{})
	Error(fields ...interface{})
	Warning(fields ...interface{})
}

type logger struct {
	Module *logrus.Logger
}

func NewLogger() *logger {
	l := &logger{}


	formatter := &logrus.JSONFormatter{}
	// formatter.SetColorScheme(&prefixed.ColorScheme{
	// 	InfoLevelStyle:  "green",
	// 	WarnLevelStyle:  "yellow",
	// 	ErrorLevelStyle: "red",
	// 	FatalLevelStyle: "red",
	// 	PanicLevelStyle: "red",
	// 	DebugLevelStyle: "blue",
	// 	PrefixStyle:     "cyan",
	// 	TimestampStyle:  "black+h",
	// 	// InfoLevelStyle: "green+b",
	// 	// PrefixStyle:    "blue+b",
	// 	// TimestampStyle: "white+h",
	// })

	logrus.ParseLevel("debug")
	logrus.SetFormatter(formatter)
	logrus.SetReportCaller(true)
	l.Module = logrus.StandardLogger()

	return l
}

func (l *logger) Info(fields ...interface{}) {
	l.Module.Info(fields...)
}

func (l *logger) Error(fields ...interface{}) {
	l.Module.Error(fields...)
}

func (l *logger) Warning(fields ...interface{}) {
	l.Module.Warn(fields...)
}

func (l *logger) Fatal(fields ...interface{}) {
	l.Module.Fatal(fields ...)
}
