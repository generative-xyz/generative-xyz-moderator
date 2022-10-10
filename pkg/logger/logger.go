package logger

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// AutoLogger is a utility struct for logging data in an extremely high performance system.
// We can use both Logger and SugarLog for logging. For more information,
// just visit https://godoc.org/go.uber.org/zap
type AutoLogger struct {
	// Sugar for logging
	*zap.SugaredLogger
	// configuration
	config map[string]interface{}
	// Logger for logging
	Logger *zap.Logger
}

func (atl *AutoLogger) Print(args ...interface{}) {
	atl.Info(args...)
}

func (atl *AutoLogger) Printf(f string, args ...interface{}) {
	atl.Infof(f, args...)
}

func (atl *AutoLogger) Println(args ...interface{}) {
	atl.Info(args)
}

// logger ddtrace.Logger
func (atl *AutoLogger) Log(msg string) {
	atl.Info(msg)
}

// Extract takes the call-scoped Logger from grpc_zap middleware.
// It always returns a Logger that has all the grpc_ctxtags updated.
func (atl *AutoLogger) Extract(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
}

// Return fields DataDog traceid
func (atl *AutoLogger) WithContext(ctx context.Context) []zapcore.Field {
	fields := []zapcore.Field{}
	span, found := tracer.SpanFromContext(ctx)
	if found {
		fields = append(fields,
			zap.Uint64("trace.traceid", span.Context().TraceID()),
			zap.Uint64("trace.spanid", span.Context().TraceID()))
	}
	return fields
}

// AtLog is logger
var AtLog *AutoLogger

func init() {
	InitLoggerDefaultDev()
}

// InitLoggerDefault -- format json
func InitLoggerDefault(enableDebug bool) {
	// init production encoder config
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.MessageKey = "message"
	// init production config
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderCfg
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	if enableDebug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	// build logger
	logger, _ := cfg.Build()

	sugarLog := logger.Sugar()
	cfgParams := make(map[string]interface{})
	AtLog = &AutoLogger{sugarLog, cfgParams, logger}
}

// InitLoggerDefaultDev -- format text
func InitLoggerDefaultDev() {
	// init development encoder config
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	// init development config
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig = encoderCfg
	cfg.OutputPaths = []string{"stdout"}
	// build logger
	logger, _ := cfg.Build()

	sugarLog := logger.Sugar()
	cfgParams := make(map[string]interface{})
	AtLog = &AutoLogger{sugarLog, cfgParams, logger}
}
