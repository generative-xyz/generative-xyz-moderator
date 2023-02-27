package logger

import "go.uber.org/zap"

type Ilogger interface {
	LogAny(message string, fields ...zap.Field)
	Info(fields ...interface{})
	Error(fields ...interface{})
	ErrorAny(message string, fields ...zap.Field)
	Warning(fields ...interface{})
	Infof(format string, fields ...interface{})
}

type logger struct {
	Module *autoLogger
}

func NewLogger(enableDebug bool) *logger {
	l := &logger{}
	log := InitLoggerDefault(enableDebug)
	l.Module = log
	return l
}

func (l *logger) Info(fields ...interface{}) {
	l.Module.SugaredLogger.Info(fields...)
}

func (l *logger) Infof(format string, fields ...interface{}) {
	l.Module.SugaredLogger.Infof(format, fields)
}

func (l *logger) Error(fields ...interface{}) {
	l.Module.SugaredLogger.Error(fields...)
}

func (l *logger) Warning(fields ...interface{}) {
	l.Module.SugaredLogger.Warn(fields...)
}

func (l *logger) Fatal(fields ...interface{}) {
	l.Module.SugaredLogger.Fatal(fields...)
}

func (l *logger) LogAny(message string, fields ...zap.Field) {
	l.Module.Logger.Info(message, fields...)
}

func (l *logger) ErrorAny(message string, fields ...zap.Field) {
	l.Module.Logger.Error(message, fields...)
}
