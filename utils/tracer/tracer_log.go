package tracer

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"rederinghub.io/utils/logger"

	"github.com/getsentry/sentry-go"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
)

type TraceLog struct {
	data         map[string]string
	tags         map[string]string
	errorKey     *string
	errorMessage *string
	error        *error
	isError      *bool
	ignoreFields []string
	log          logger.Ilogger
	tracerID     *string
}

func NewTraceLog() *TraceLog {
	t := new(TraceLog)
	t.data = make(map[string]string)
	t.tags = make(map[string]string)
	l := logger.NewLogger()

	t.log = l
	t.ignoreFields = []string{"diagram"}
	return t
}

func (t *TraceLog) SetData(key string, value interface{}) {
	parsedValue, _ := json.Marshal(value)
	t.data[key] = string(parsedValue)

	if t.log != nil {
		f := zap.Any(key, value)
		t.log.Info(key, f)
	}
}

func (t *TraceLog) SetTag(key string, value interface{}) {
	parsedValue, _ := json.Marshal(value)
	t.tags[key] = string(parsedValue)
}

func (t *TraceLog) Error(key string, value string, err error) {
	isErrr := true
	t.isError = &isErrr
	t.errorKey = &key
	t.errorMessage = &value
	t.error = &err

	if t.log != nil {
		t.log.Error(key, err)
	}
}

func (t TraceLog) GetData() map[string]string {
	return t.data
}

func (t *TraceLog) ToSpan(span opentracing.Span) {
	tracerID := GetTraceID(span)
	t.tracerID = &tracerID

	for key, tag := range t.tags {
		tag = strings.ReplaceAll(tag, `"`, "")
		tag = strings.TrimSpace(tag)
		span.SetTag(key, tag)
	}

	for key, log := range t.data {
		span.LogFields(
			LogString(key, log),
		)
	}

	if t.isError != nil {
		if *t.isError {
			span.SetTag("error", true)
			errMess := fmt.Sprintf("[%s][%s][%s]", *t.tracerID, *t.errorKey, *t.errorMessage)
			sentry.CaptureException(errors.Wrap(*t.error, errMess))
			SpanError(span, *t.error, *t.errorKey, *t.errorMessage)
		}

	}
}

func GetTraceID(span opentracing.Span) string {

	tracerID := ""
	if jaegerCtx, ok := span.Context().(jaeger.SpanContext); ok {
		tID := jaegerCtx.TraceID().String()
		//sID := jaegerCtx.SpanID().String()
		tracerID = tID
	}

	return tracerID
}
