package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type IMiddleware interface {
	LoggingMiddleware(next http.Handler) http.Handler
	AccessToken(next http.Handler) http.Handler
	AccessTokenPassThrough(next http.Handler) http.Handler
	AuthorizationFunc(next http.Handler) http.Handler
}

type middleware struct {
	log              logger.Ilogger
	usecase          usecase.Usecase
	response         response.IHttpResponse
	cache            redis.IRedisCache
	cacheAuthService redis.IRedisCache
}

func NewMiddleware(uc usecase.Usecase, g *global.Global) *middleware {
	m := new(middleware)
	m.log = g.Logger
	m.usecase = uc
	m.response = response.NewHttpResponse()
	m.cache = g.Cache
	m.cacheAuthService = g.CacheAuthService
	return m
}

func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// defer func() {
		// 	if err := recover(); err != nil {
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		m.log.Error(
		// 			"err", err,
		// 			"trace", debug.Stack(),
		// 		)
		// 	}
		// }()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		logger.AtLog.Info(fmt.Sprintf("Request:[%s] %s - status: %d - duration %s =====", r.Method, r.URL.EscapedPath(), wrapped.status, time.Since(start)))
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

// Authenticate
func (m *middleware) AccessToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get(utils.AUTH_TOKEN)
		if token == "" {
			err := errors.New("token is empty")
			logger.AtLog.Logger.Error("token_is_empty", zap.Error(err))
			m.response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
			return
		}

		token = helpers.ReplaceToken(token)

		//TODO implement here
		p, err := m.usecase.ValidateAccessToken(token)
		if err != nil {
			logger.AtLog.Logger.Error("cannot_verify_token", zap.Error(err))
			m.response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
			return
		}

		logger.AtLog.Logger.Info("AccessToken", zap.Any("profile", p))
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

// Authenticate
func (m *middleware) AccessTokenPassThrough(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get(utils.AUTH_TOKEN)
		if token == "" {
			err := errors.New("token is empty")
			logger.AtLog.Logger.Error("token_is_empty", zap.Error(err))
			next.ServeHTTP(w, r.WithContext(r.Context()))
			return
		}

		token = helpers.ReplaceToken(token)

		//TODO implement here
		p, err := m.usecase.ValidateAccessToken(token)
		if err != nil {
			logger.AtLog.Logger.Error("cannot_verify_token", zap.Error(err))
			next.ServeHTTP(w, r.WithContext(r.Context()))
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

// Just set SIGNED_WALLET_ADDRESS, SIGNED_USER_ID to context if have
// Authorization
func (m *middleware) AuthorizationFunc(next http.Handler) http.Handler {
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
		ctx = context.WithValue(ctx, utils.AUTH_TOKEN, token)
		ctx = context.WithValue(ctx, utils.SIGNED_WALLET_ADDRESS, p.WalletAddress)
		ctx = context.WithValue(ctx, utils.SIGNED_USER_ID, p.Uid)
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
