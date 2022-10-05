package middleware

import (
	"context"
	"google.golang.org/grpc"
	"rederinghub.io/pkg/log"
	"time"
)

type Interceptor struct {
	logger log.Logger
}

func NewInterceptor(logger log.Logger) Interceptor {
	return Interceptor{
		logger: logger,
	}
}

func (i Interceptor) WithTimeoutInterceptor() grpc.UnaryServerInterceptor {
	return i.contextInterceptor
}

func (i Interceptor) contextInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return handler(timeoutCtx, req)
}
