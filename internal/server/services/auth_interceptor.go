package services

import (
	"context"

	"github.com/OvsienkoValeriya/GophKeeper/internal/logger"
	"github.com/OvsienkoValeriya/GophKeeper/internal/server/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtConfig *auth.JWTConfig
}

type ContextKey string

const UserIDKey ContextKey = "userID"

func NewAuthInterceptor(jwtConfig *auth.JWTConfig) *AuthInterceptor {
	return &AuthInterceptor{
		jwtConfig: jwtConfig,
	}
}

func (interceptor *AuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Sugar.Info("UnaryInterceptor: ", info.FullMethod)
		newCtx, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func (interceptor *AuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger.Sugar.Info("StreamInterceptor: ", info.FullMethod)
		_, err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		return handler(srv, stream)
	}
}

var publicMethods = map[string]bool{
	"/gophkeeper.auth.AuthService/Login":    true,
	"/gophkeeper.auth.AuthService/Register": true,
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, fullMethod string) (context.Context, error) {
	if publicMethods[fullMethod] {
		return ctx, nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "token is not provided")
	}
	accessToken := token[0]
	claims, err := interceptor.jwtConfig.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	newCtx := context.WithValue(ctx, UserIDKey, claims.UserID)
	return newCtx, nil
}
