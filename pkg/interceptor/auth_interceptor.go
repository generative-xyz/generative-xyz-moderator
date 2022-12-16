package interceptor

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"rederinghub.io/pkg/oauth2service"
)

type ContextKey string

const (
	ContextKeyToken  ContextKey = "token"
	ContextKeyuserID ContextKey = "id"
	AuthKey string = "authorization"
)

type AuthInterceptor struct {
	allowedAuthMethods []string
	auth oauth2service.Auth2
}

func NewAuthInterceptor(auth oauth2service.Auth2, AllowedAuthMethods []string) *AuthInterceptor {
	return &AuthInterceptor{
		allowedAuthMethods: AllowedAuthMethods,
		auth: auth,
	}
}

func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		isCheckedMiddleWare := false
		method := ai.extractMethod(info.FullMethod)
		for _, m := range ai.allowedAuthMethods {
			if method == m {
				isCheckedMiddleWare = true
			}
		}

		if isCheckedMiddleWare {
			ctx, err := ai.authorize(ctx)
			if err != nil {
				return nil, status.New(codes.Internal, err.Error()).Err()
			}
			return handler(ctx, req)
		}
		return handler(ctx, req)
	}
}

func (ai *AuthInterceptor) authorize(ctx context.Context) (context.Context, error) {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(m[AuthKey]) == 0 {
		return ctx, status.New(codes.Unauthenticated, "missing token").Err()
	}

	authData, err := ai.auth.ClaimToken(m[AuthKey][0])
	if err != nil {
		return ctx, status.New(codes.Unauthenticated, "unauthorized").Err()
	}

	var meta map[string]interface{}
	b, err := json.Marshal(authData)
	if err != nil {
		return ctx, status.New(codes.Unauthenticated, "unauthorized").Err()
	} else {
		if err := json.Unmarshal(b, &meta); err != nil {
			log.Println("Error while unmarshaling authData data", err)
		}
	}
	meta[string(ContextKeyToken)] = m[AuthKey][0]
	spew.Dump(meta)
	ctx = ai.SetContextMetadata(ctx, meta)
	return ctx, nil
}

func (ai *AuthInterceptor) extractMethod(fullMethod string) string {
	re := regexp.MustCompile(`.+/(\w+)$`)
	method := re.ReplaceAllString(fullMethod, "$1")

	return method
}

func (ai *AuthInterceptor) SetContextMetadata(ctx context.Context, data map[string]interface{}) context.Context {
	for k, v := range data {
	    key := strings.ToLower(k)
		ctx = context.WithValue(ctx, ContextKey(key), v)
	}

	return ctx
}
