package log

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"io"
	"os"
	"rederinghub.io/pkg/config"
	"strings"
)

var AppLogger = NewLogger("app_default")

type Logger interface {
	Output(w io.Writer) zerolog.Logger

	With() zerolog.Context

	Level(level zerolog.Level) zerolog.Logger

	Sample(s zerolog.Sampler) zerolog.Logger

	Hook(h zerolog.Hook) zerolog.Logger

	Err(err error) *zerolog.Event

	Trace() *zerolog.Event

	Debug() *zerolog.Event

	Info() *zerolog.Event

	Warn() *zerolog.Event

	Error() *zerolog.Event

	Fatal() *zerolog.Event

	Panic() *zerolog.Event

	WithLevel(level zerolog.Level) *zerolog.Event

	Prefix(prefix string) zerolog.Context

	CtxPrefix(ctx context.Context, prefix string) zerolog.Context

	Log() *zerolog.Event

	Print(v ...interface{})

	Printf(format string, v ...interface{})

	Ctx(ctx context.Context) *zerolog.Logger
}

func NewLogger(srv string) Logger {
	logger := zerolog.New(os.Stderr).
		With().
		Caller().
		Str("service", strings.ReplaceAll(srv, " ", "_")).
		Timestamp().
		Logger()
	logger = logger.Level(GetLogLevel(config.ServerConfig().LogLevel))
	return &ZeroLog{
		logger: logger,
	}
}

func NewZeroLog(logger zerolog.Logger) Logger {
	return &ZeroLog{
		logger: logger,
	}
}

type ZeroLog struct {
	logger zerolog.Logger
}

func (l *ZeroLog) Output(w io.Writer) zerolog.Logger {
	return l.logger.Output(w)
}

func (l *ZeroLog) With() zerolog.Context {
	return l.logger.With()
}

func (l *ZeroLog) Level(level zerolog.Level) zerolog.Logger {
	return l.logger.Level(level)
}

func (l *ZeroLog) Sample(s zerolog.Sampler) zerolog.Logger {
	return l.logger.Sample(s)
}

func (l *ZeroLog) Hook(h zerolog.Hook) zerolog.Logger {
	return l.logger.Hook(h)
}

func (l *ZeroLog) Err(err error) *zerolog.Event {
	return l.logger.Err(err)
}

func (l *ZeroLog) Trace() *zerolog.Event {
	return l.logger.Trace()
}

func (l *ZeroLog) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *ZeroLog) Info() *zerolog.Event {
	return l.logger.Info()
}

func (l *ZeroLog) Warn() *zerolog.Event {
	return l.logger.Warn()
}

func (l *ZeroLog) Error() *zerolog.Event {
	return l.logger.Error()
}

func (l *ZeroLog) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

func (l *ZeroLog) Panic() *zerolog.Event {
	return l.logger.Panic()
}

func (l *ZeroLog) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

func (l *ZeroLog) Log() *zerolog.Event {
	return l.logger.Log()
}

func (l *ZeroLog) Print(v ...interface{}) {
	l.logger.Print(v...)
}

func (l *ZeroLog) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *ZeroLog) Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func (l *ZeroLog) Prefix(prf string) zerolog.Context {
	return l.With().Str("prefix", prf)
}

func (l *ZeroLog) CtxPrefix(ctx context.Context, prf string) zerolog.Context {
	return l.With().Str("request-id", fmt.Sprint(ctx.Value(middleware.RequestIDKey))).Str("prefix", prf)
}

func GetLogLevel(level string) zerolog.Level {
	switch level {
	case "info":
		return zerolog.InfoLevel
	case "error":
		return zerolog.ErrorLevel
	case "warning":
		return zerolog.WarnLevel
	case "disabled":
		return zerolog.Disabled
	}
	return zerolog.DebugLevel
}
