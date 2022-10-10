package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Interceptor struct {
}

func NewInterceptor() Interceptor {
	return Interceptor{}
}

func (i Interceptor) WithTimeoutInterceptor() grpc.UnaryServerInterceptor {
	return i.contextInterceptor
}

func (i Interceptor) contextInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return handler(timeoutCtx, req)
}
