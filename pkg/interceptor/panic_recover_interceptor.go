package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type PanicRecoveryInterceptor struct {
}

func NewPanicRecoveryInterceptor() *PanicRecoveryInterceptor {
	return &PanicRecoveryInterceptor{}
}

func (pri *PanicRecoveryInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered from err: ", err)
			}
		}()

		return handler(ctx, req)
	}
}
