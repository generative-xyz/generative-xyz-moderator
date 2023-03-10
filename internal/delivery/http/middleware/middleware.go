package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"rederinghub.io/internal/delivery/http/response"
	_httpResponse "rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type IMiddleware interface {
	LoggingMiddleware(next http.Handler) http.Handler
	VerifyEmail(next http.Handler) http.Handler
	AccessToken(next http.Handler) http.Handler
	UserToken(next http.Handler) http.Handler
	AuthorizeFunc(next http.Handler) http.Handler
}

type middleware struct {
	log              logger.Ilogger
	usecase          usecase.Usecase
	response         _httpResponse.IHttpResponse
	cache            redis.IRedisCache
	cacheAuthService redis.IRedisCache
}

func NewMiddleware(uc usecase.Usecase, g *global.Global) *middleware {
	m := new(middleware)
	m.log = g.Logger
	m.usecase = uc
	m.response = _httpResponse.NewHttpResponse()
	m.cache = g.Cache
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

func (m *middleware) AccessToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get(utils.AUTH_TOKEN)
		if token == "" {
			err := errors.New("token is empty")
			m.log.Error("token_is_empty", "token is empty", err)
			m.response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
			return
		}

		token = helpers.ReplaceToken(token)

		//TODO implement here
		p, err := m.usecase.ValidateAccessToken(token)
		if err != nil {
			m.log.Error("cannot_verify_token", "token cannot be verified", err)
			m.response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
			return
		}

		m.log.Info("profile", p)
		m.cache.SetData(helpers.GenerateCachedProfileKey(token), p)
		m.cache.SetStringData(helpers.GenerateUserKey(token), p.Uid)

		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.AUTH_TOKEN, token)
		ctx = context.WithValue(ctx, utils.SIGNED_WALLET_ADDRESS, p.WalletAddress)
		//ctx = context.WithValue(ctx, utils.SIGNED_EMAIL, p.Email)
		ctx = context.WithValue(ctx, utils.SIGNED_USER_ID, p.Uid)
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (m *middleware) UserToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		token := r.Header.Get(utils.AUTH_TOKEN)
		if token != "" {
			token = helpers.ReplaceToken(token)
			p, err := m.usecase.ValidateAccessToken(token)
			if err == nil {
				m.log.Info("profile", p)
				//m.log.Info("signedDetail", signedDetail)
				m.cache.SetData(helpers.GenerateCachedProfileKey(token), p)
				m.cache.SetStringData(helpers.GenerateUserKey(token), p.Uid)

				ctx = context.WithValue(ctx, utils.AUTH_TOKEN, token)
				ctx = context.WithValue(ctx, utils.SIGNED_WALLET_ADDRESS, p.WalletAddress)
				//ctx = context.WithValue(ctx, utils.SIGNED_EMAIL, p.Email)
				ctx = context.WithValue(ctx, utils.SIGNED_USER_ID, p.Uid)
			}

		}

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

// Just set SIGNED_WALLET_ADDRESS, SIGNED_USER_ID to context if have
// Authorize not Authenticate
func (m *middleware) AuthorizeFunc(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := helpers.ReplaceToken(r.Header.Get(utils.AUTH_TOKEN))
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
		p, err := m.usecase.ValidateAccessToken(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.SIGNED_WALLET_ADDRESS, p.WalletAddress)
		ctx = context.WithValue(ctx, utils.SIGNED_USER_ID, p.Uid)
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
