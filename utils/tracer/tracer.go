package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"rederinghub.io/utils/logger"

	"github.com/davecgh/go-spew/spew"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

type ITracer interface {
	StartSpan(operationName string) opentracing.Span
	StartSpanFromInjection(tracingInjection map[string]string, operationName string) opentracing.Span
	StartSpanFromHeaderInjection(tracingInjection http.Header, operationName string) opentracing.Span
	SpanError(span opentracing.Span, err error, messageKey string, message string) (interface{}, error)
	LogObject(key string, value interface{}) log.Field
	LogString(key string, value string) log.Field
	LogInt(key string, value int) log.Field
	StartSpanFromRoot(rootspan opentracing.Span, optName string) opentracing.Span
	StartWithOpts(optName string, opts ...opentracing.StartSpanOption) opentracing.Span
	GetTrace() opentracing.Tracer
	CtxWithSpan(ctx context.Context, span opentracing.Span) context.Context
	StartSpanFromContext(ctx context.Context, name string) opentracing.Span
	FinishSpan(span opentracing.Span, log *TraceLog)
}

type tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
	logger logger.Ilogger
}

func NewTracing(log logger.Ilogger) *tracer {
	t := new(tracer)
	cfg, err := config.FromEnv()
	if err != nil {
		//panic(fmt.Sprintf("Could not parse Jaeger env vars: %s", err.Error()))
	}

	cfg.Sampler = &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	t.closer = closer
	t.tracer = tracer
	t.logger = log
	return t
}

func (t *tracer) StartSpan(operationName string) opentracing.Span {
	return t.tracer.StartSpan(operationName)
}

func (t *tracer) StartWithOpts(optName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.tracer.StartSpan(optName, opts...)
}

func (t *tracer) StartSpanFromRoot(rootspan opentracing.Span, optName string) opentracing.Span {
	span := t.tracer.StartSpan(optName, opentracing.ChildOf(rootspan.Context()))
	return span
}

func (t *tracer) StartSpanFromInjection(tracingInjection map[string]string, operationName string) opentracing.Span {
	tracer := t.tracer
	spanCtx, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(tracingInjection))

	if err != nil {
		spew.Dump(err)
		span := tracer.StartSpan(operationName)
		return span
	}

	span := tracer.StartSpan(operationName, ext.RPCServerOption(spanCtx))
	return span
}

func (t *tracer) StartSpanFromHeaderInjection(tracingInjection http.Header, operationName string) opentracing.Span {
	tracer := t.tracer
	spanCtx, err := t.GetTrace().Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(tracingInjection))
	if err != nil {
		spew.Dump(err)
		span := tracer.StartSpan(operationName)
		return span
	}

	span := t.StartWithOpts(operationName, ext.RPCServerOption(spanCtx))
	return span
}

func (t *tracer) SpanError(span opentracing.Span, err error, messageKey string, message string) (interface{}, error) {
	span.LogFields(
		log.Error(err),
		log.String(messageKey, message),
	)
	ext.Error.Set(span, true)

	t.logger.Info("SpanError", zap.String(messageKey, message))
	t.logger.Error("SpanError", zap.Error(err))

	return nil, err
}

func (t *tracer) LogObject(key string, value interface{}) log.Field {
	bytesArray, err := json.Marshal(value)
	if err != nil {
		return log.String(key, "")
	}

	t.logger.Info("LogObject", zap.Any(key, value))

	jsonMessage := string(bytesArray)
	return log.Object(key, jsonMessage)
}

func (t *tracer) LogString(key string, value string) log.Field {
	t.logger.Info("LogString", zap.String(key, value))
	return log.Object(key, value)
}

func (t *tracer) LogInt(key string, value int) log.Field {
	t.logger.Info("LogInt", zap.Int(key, value))
	return log.Int(key, value)
}

func (t *tracer) GetTrace() opentracing.Tracer {
	return t.tracer
}

func SpanError(span opentracing.Span, err error, messageKey string, message string) (interface{}, error) {
	span.LogFields(
		log.Error(err),
		log.String(messageKey, message),
	)
	//sentry.CaptureException(err)
	ext.Error.Set(span, true)
	return nil, err
}

func LogString(key string, value string) log.Field {
	return log.String(key, value)
}

func (t *tracer) CtxWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(ctx, span)
}

func (t *tracer) StartSpanFromContext(ctx context.Context, name string) opentracing.Span {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		return t.tracer.StartSpan(name)
	}

	return t.tracer.StartSpan(name, opentracing.ChildOf(parentSpan.Context()))
}

func (t *tracer) FinishSpan(span opentracing.Span, log *TraceLog) {

	log.ToSpan(span)

	span.Finish()
}
