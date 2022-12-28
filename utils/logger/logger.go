package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Ilogger interface {
	Info(message string, fields ...zapcore.Field)
	Error(message string, fields ...interface{})
	Warning(message string, fields ...interface{})
}

type logger struct {
	Module *zap.Logger
}

func NewLogger() *logger {
	log := new(logger)

	m, err := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message", // <--
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		DisableCaller: true,
	}.Build()
	if err != nil {
		panic(err)
	}
	log.Module = m
	return log
}

func (l *logger) Info(message string, fields ...zapcore.Field) {
	l.Module.Info(message, fields...)
}

func (l *logger) Error(message string, fields ...interface{}) {
	l.Module.Sugar().Error(message, fields)
}

func (l *logger) Warning(message string, fields ...interface{}) {
	l.Module.Sugar().Warn(message, fields)
}

func (l *logger) Fatal(message string, fields ...interface{}) {
	l.Module.Sugar().Fatal(message, fields)
}
