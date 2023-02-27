package logger

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// AutoLogger is a utility struct for logging data in an extremely high performance system.
// We can use both Logger and SugarLog for logging. For more information,
// just visit https://godoc.org/go.uber.org/zap
type autoLogger struct {
	// Sugar for logging
	*zap.SugaredLogger
	// configuration
	config map[string]interface{}
	// Logger for logging
	Logger *zap.Logger
}

func (atl *autoLogger) Print(args ...interface{}) {
	atl.Info(args...)
}

func (atl *autoLogger) Printf(f string, args ...interface{}) {
	atl.Infof(f, args...)
}

func (atl *autoLogger) Println(args ...interface{}) {
	atl.Info(args)
}

// logger ddtrace.Logger
func (atl *autoLogger) Log(msg string) {
	atl.Info(msg)
}

// Extract takes the call-scoped Logger from grpc_zap middleware.
// It always returns a Logger that has all the grpc_ctxtags updated.
func (atl *autoLogger) Extract(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
	//return nil
}

// Return fields DataDog traceid
func (atl *autoLogger) WithContext(ctx context.Context) []zapcore.Field {
	fields := []zapcore.Field{}
	span, found := tracer.SpanFromContext(ctx)
	if found {
		fields = append(fields,
			zap.Uint64("trace.traceid", span.Context().TraceID()),
			zap.Uint64("trace.spanid", span.Context().SpanID()))
	}
	return fields
}

// LogRoundTrip prints the information about request and response.
func (atl *autoLogger) LogRoundTrip(
	req *http.Request,
	res *http.Response,
	err error,
	start time.Time,
	dur time.Duration,
) error {
	var (
		nReq int64
		nRes int64
	)

	// Count number of bytes in request and response.
	if req != nil && req.Body != nil && req.Body != http.NoBody {
		nReq, _ = io.Copy(ioutil.Discard, req.Body)
	}
	if res != nil && res.Body != nil && res.Body != http.NoBody {
		nRes, _ = io.Copy(ioutil.Discard, res.Body)
	}

	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.Int("status_code", res.StatusCode),
		zap.Duration("duration", dur),
		zap.Int64("req_bytes", nReq),
		zap.Int64("res_bytes", nRes),
	}

	// Set error level.
	switch {
	case err != nil:
		atl.Logger.With(fields...).Error(req.URL.String())
	case res != nil && res.StatusCode > 0 && res.StatusCode < 300:
		atl.Logger.With(fields...).Info(req.URL.String())
	case res != nil && res.StatusCode > 299 && res.StatusCode < 500:
		atl.Logger.With(fields...).Warn(req.URL.String())
	case res != nil && res.StatusCode > 499:
		atl.Logger.With(fields...).Error(req.URL.String())
	default:
		atl.Logger.With(fields...).Error(req.URL.String())
	}

	return nil
}

// RequestBodyEnabled makes the client pass request body to logger
func (atl *autoLogger) RequestBodyEnabled() bool { return true }

// RequestBodyEnabled makes the client pass response body to logger
func (atl *autoLogger) ResponseBodyEnabled() bool { return true }

// AtLog is logger
var atLog *autoLogger

func init() {
	InitLoggerDefaultDev()
}

// InitLoggerDefault -- format json
func InitLoggerDefault(enableDebug bool) *autoLogger {
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
	atLog = &autoLogger{sugarLog, cfgParams, logger}
	return atLog
}

// InitLoggerDefaultDev -- format text
func InitLoggerDefaultDev() *autoLogger {
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
	atLog = &autoLogger{sugarLog, cfgParams, logger}
	return atLog
}

// GetLoggerInstanceFromContext returns the logger instance from context
func GetLoggerInstanceFromContext(ctx context.Context) *zap.Logger {
	return atLog.Logger.With(atLog.WithContext(ctx)...)
}