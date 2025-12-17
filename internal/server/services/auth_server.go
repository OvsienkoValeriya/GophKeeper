package services

import (
	"context"
	"fmt"

	pb "github.com/OvsienkoValeriya/GophKeeper/api/gen"
	"github.com/OvsienkoValeriya/GophKeeper/internal/models"
	"github.com/OvsienkoValeriya/GophKeeper/internal/server/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	userStore UserStore
	jwtConfig *auth.JWTConfig
}

func NewAuthServer(userStore UserStore, jwtConfig *auth.JWTConfig) *AuthServer {
	return &AuthServer{
		userStore: userStore,
		jwtConfig: jwtConfig,
	}
}

func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.userStore.GetUserByUsername(ctx, req.GetUsername())

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	if !auth.ValidatePassword(user.Password, req.GetPassword()) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	accessToken, err := server.jwtConfig.GenerateJWT(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	return &pb.LoginResponse{
		UserId:       proto.String(fmt.Sprintf("%d", user.ID)),
		AccessToken:  proto.String(accessToken),
		RefreshToken: proto.String(""),
	}, nil
}

func (server *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	if username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username is required")
	}
	if len(password) < 8 {
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 8 characters")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	user := &models.User{
		Username: username,
		Password: hashedPassword,
	}

	createdUser, err := server.userStore.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	accessToken, err := server.jwtConfig.GenerateJWT(createdUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token")
	}

	return &pb.RegisterResponse{
		UserId:       proto.String(fmt.Sprintf("%d", createdUser.ID)),
		AccessToken:  proto.String(accessToken),
		RefreshToken: proto.String(""),
	}, nil
}

func (server *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	user, err := server.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user")
	}

	newAccessToken, err := server.jwtConfig.GenerateJWT(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token")
	}

	newRefreshToken, err := server.jwtConfig.GenerateRefreshToken(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token")
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  proto.String(newAccessToken),
		RefreshToken: proto.String(newRefreshToken),
	}, nil
}

func (server *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {

	return &pb.LogoutResponse{Success: proto.Bool(true)}, nil
}

func (server *AuthServer) SetMasterKey(ctx context.Context, req *pb.SetMasterKeyRequest) (*pb.SetMasterKeyResponse, error) {

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	hasMasterKey, err := server.userStore.HasMasterKey(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check master key: %v", err)
	}
	if hasMasterKey {
		return nil, status.Error(codes.AlreadyExists, "master key already set")
	}

	if err := server.userStore.SetMasterKey(ctx, userID, req.GetSalt(), req.GetVerifier()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set master key: %v", err)
	}

	return &pb.SetMasterKeyResponse{Success: proto.Bool(true)}, nil
}

func (server *AuthServer) GetMasterKeyData(ctx context.Context, req *pb.GetMasterKeyDataRequest) (*pb.GetMasterKeyDataResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	salt, verifier, err := server.userStore.GetMasterKeyData(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get master key data: %v", err)
	}

	hasMasterKey := len(salt) > 0 && len(verifier) > 0

	return &pb.GetMasterKeyDataResponse{
		Salt:         salt,
		Verifier:     verifier,
		HasMasterKey: proto.Bool(hasMasterKey),
	}, nil
}

func (server *AuthServer) HasMasterKey(ctx context.Context, req *pb.HasMasterKeyRequest) (*pb.HasMasterKeyResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	hasMasterKey, err := server.userStore.HasMasterKey(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check master key: %v", err)
	}

	return &pb.HasMasterKeyResponse{HasMasterKey: proto.Bool(hasMasterKey)}, nil
}
