package client

import (
	"context"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	tokenStore  *FileTokenStore
	mu          sync.RWMutex
}

func NewAuthInterceptor(authClient *AuthClient, authMethods map[string]bool, tokenStore *FileTokenStore) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
		tokenStore:  tokenStore,
	}
	return interceptor, nil
}

func (interceptor *AuthInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("UnaryInterceptor: ", method)
		if interceptor.authMethods[method] {
			ctx = interceptor.attachToken(ctx)
		}
		return invoker(ctx, method, req, reply, cc, opts...)

	}

}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	accessToken, _, err := interceptor.tokenStore.LoadTokens()
	if err != nil {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)
}

func (interceptor *AuthInterceptor) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Println("StreamInterceptor: ", method)
		if interceptor.authMethods[method] {
			ctx = interceptor.attachToken(ctx)
		}
		return streamer(ctx, desc, cc, method, opts...)

	}

}
