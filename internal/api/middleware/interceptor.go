package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"rederinghub.io/pkg/interceptor"
	"rederinghub.io/pkg/oauth2service"
)

var AllowedAuthMethods = []string{
	"GetProfile",
	"UpdateProfile",
}


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

func (i Interceptor) AuthInterceptor() grpc.UnaryServerInterceptor {	
	auth := oauth2service.NewAuth2()
	authInterceptor := interceptor.NewAuthInterceptor(*auth, AllowedAuthMethods)
	return authInterceptor.Unary()
}
