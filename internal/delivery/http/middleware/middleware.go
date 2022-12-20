package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"rederinghub.io/internal/delivery/http/response"
	_httpResponse "rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type IMiddleware interface {
	LoggingMiddleware(next http.Handler) http.Handler
	VerifyEmail(next http.Handler) http.Handler
	AccessToken(next http.Handler) http.Handler
	Tracer(next http.Handler) http.Handler
}

type middleware struct {
	log           logger.Ilogger
	usecase       usecase.Usecase
	response      _httpResponse.IHttpResponse
	cache         redis.IRedisCache
	cacheAuthService         redis.IRedisCache
	tracer        tracer.ITracer
}

func NewMiddleware(uc usecase.Usecase, g *global.Global) *middleware {
	m := new(middleware)
	m.log = g.Logger
	m.usecase = uc
	m.response = _httpResponse.NewHttpResponse()
	m.cache = g.Cache
	m.tracer = g.Tracer
	m.cacheAuthService = g.CacheAuthService
	return m
}

func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				m.log.Error(
					"err", err,
					"trace", debug.Stack(),
				)
			}
		}()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		m.log.Info(fmt.Sprintf("Request:[%s] %s - status: %d - duration %s =====", r.Method, r.URL.EscapedPath(), wrapped.status, time.Since(start)))
	}

	return http.HandlerFunc(fn)
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (m *middleware) Tracer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := m.tracer.GetTrace().Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(r.Header))
		method := r.Method
		path := r.URL.Path

		vars := mux.Vars(r)
		for key, value := range vars {
			path = strings.ReplaceAll(path, value, key)
		}

		spanName := fmt.Sprintf("[%s] - %s", method, path)
		span := m.tracer.StartWithOpts(spanName, ext.RPCServerOption(spanCtx))
		log := tracer.NewTraceLog()

		defer log.ToSpan(span)
		defer span.Finish()

		requestURI := r.RequestURI
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		dataIn := bytes.NewBuffer(body)
		r.Body = ioutil.NopCloser(dataIn)

		bodyData := make(map[string]interface{})
		json.Unmarshal(body, &bodyData)

		log.SetData("method", method)
		log.SetData("requestURL", requestURI)
		log.SetData("header", r.Header)
		log.SetData("body", bodyData)

		
		m.tracer.GetTrace().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
	}

	return http.HandlerFunc(fn)
}

func (m *middleware) AccessToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := m.tracer.GetTrace().Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(r.Header))

		//requestURI := r.RequestURI
		method := r.Method
		spanName := fmt.Sprintf("[%s][%s] - %s ","AccessToken", method, r.URL.Path)
		
		span := m.tracer.StartWithOpts(spanName, ext.RPCServerOption(spanCtx))
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)


		token := r.Header.Get(utils.AUTH_TOKEN)
		if token == "" {
			err := errors.New("token is empty")
			log.Error("token_is_empty", "token is empty", err)
			m.response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
			return
		}

		//TODO implement here
		p, err := m.usecase.UserProfile(span, token)
		if err != nil {
			log.Error("cannot_verify_token", "token cannot be verified", err)
			m.response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
			return
		}

		log.SetData("profile", p)
		//log.SetData("signedDetail", signedDetail)
		m.tracer.GetTrace().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		log.SetTag(utils.WALLET_ADDRESS_TAG, p.WalletAddress)
		
 
		m.cache.SetData( helpers.GenerateCachedProfileKey(token), p)
		m.cache.SetStringData(helpers.GenerateUserKey(token), p.ID)
		log.SetTag(utils.EMAIL_TAG, p.Email)

		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.AUTH_TOKEN, token)
		ctx = context.WithValue(ctx, utils.SIGNED_WALLET_ADDRESS, p.WalletAddress)
		ctx = context.WithValue(ctx, utils.SIGNED_EMAIL, p.Email)
		ctx = context.WithValue(ctx, utils.SIGNED_USER_ID, p.ID)
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

//TODO - Is the inserted email belong to the correct ID.
func (m *middleware) VerifyEmail(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
	}

	return http.HandlerFunc(fn)
}