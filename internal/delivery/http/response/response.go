package response

import (
	"encoding/json"
	"net/http"

	"rederinghub.io/utils/tracer"

	"github.com/opentracing/opentracing-go"
)

type IHttpResponse interface {
	RespondWithError(w http.ResponseWriter, httpCode int, appCode int, payload error)
	RespondSuccess(w http.ResponseWriter, httpCode int, appCode int, payload interface{}, customerMessage string)
	RespondWithoutContainer(w http.ResponseWriter, httpCode int, payload interface{})
	SetTrace(tracer tracer.ITracer)
	SetSpan(rootSpan opentracing.Span)
	SetLog(tracer tracer.ITracer, rootSpan opentracing.Span)
}

type JsonResponse struct {
	Error *RespondErr      `json:"error"`
	Status    bool        `json:"status"`
	Data    interface{} `json:"data"`
}

type RespondErr struct {
	Message string `json:"message"`
	ErrorCode int `json:"code"`
}

type httpResponse struct {
	RootSpan opentracing.Span
	Tracer   tracer.ITracer
}

func NewHttpResponse() *httpResponse {
	return new(httpResponse)
}

func (h *httpResponse) SetTrace(tracer tracer.ITracer) {
	h.Tracer = tracer
}

func (h *httpResponse) SetSpan(rootSpan opentracing.Span) {
	h.RootSpan = rootSpan
}

func (h *httpResponse) SetLog(tracer tracer.ITracer, rootSpan opentracing.Span) {
	h.Tracer = tracer
	h.RootSpan = rootSpan
}

func (h *httpResponse) RespondWithError(w http.ResponseWriter, httpCode int, appCode int, payload error) {
	h.respondWithJSON(w, payload, httpCode, appCode, payload, "")
}

func (h *httpResponse) RespondSuccess(w http.ResponseWriter, httpCode int, appCode int, payload interface{}, customerMessage string) {
	h.respondWithJSON(w, nil, httpCode, appCode, payload, customerMessage)
}

func (h *httpResponse) respondWithJSON(w http.ResponseWriter, respErr error, httpCode int, appCode int, payload interface{}, customerMessage string) {

	code := ResponseMessage[appCode].Code
	//message := ResponseMessage[appCode].Message

	if customerMessage != "" {
		//message = customerMessage
	}

	jsr := JsonResponse{
		Data:    payload,
		Status:    true,
	}

	if respErr != nil {
		errMessage := &RespondErr{}
		errMessage.Message = respErr.Error()
		errMessage.ErrorCode = code
		jsr.Error = errMessage
		jsr.Status = false
	}

	if h.Tracer != nil && h.RootSpan != nil {
		span := h.Tracer.StartSpanFromRoot(h.RootSpan, "respondWithJSON")
		defer span.Finish()

		span.LogFields(
			h.Tracer.LogObject("response.Json", jsr),
			h.Tracer.LogInt("response.code", code),
		)
	}

	response, _ := json.Marshal(jsr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, err := w.Write(response)
	if err != nil {
		panic(err)
	}
}


func (h *httpResponse) RespondWithoutContainer(w http.ResponseWriter, httpCode int, payload interface{}) {

	if h.Tracer != nil && h.RootSpan != nil {
		span := h.Tracer.StartSpanFromRoot(h.RootSpan, "respondWithJSON")
		defer span.Finish()

		span.LogFields(
			h.Tracer.LogObject("response.Json", payload),
		)
	}

	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, err := w.Write(response)
	if err != nil {
		panic(err)
	}
}
