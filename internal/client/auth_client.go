package client

import (
	"context"
	"time"

	pb "github.com/OvsienkoValeriya/GophKeeper/api/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type AuthClient struct {
	service pb.AuthServiceClient
}

func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{
		service: service,
	}
}

func (client *AuthClient) Login(username, password string) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: proto.String(username),
		Password: proto.String(password),
	}

	return client.service.Login(ctx, req)
}

func (client *AuthClient) Register(username, password string) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.RegisterRequest{
		Username: proto.String(username),
		Password: proto.String(password),
	}

	return client.service.Register(ctx, req)
}

func (client *AuthClient) RefreshToken(accessToken string) (*pb.RefreshTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.RefreshTokenRequest{}

	return client.service.RefreshToken(ctx, req)
}

func (client *AuthClient) Logout(accessToken string) (*pb.LogoutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.LogoutRequest{
		AccessToken: proto.String(accessToken),
	}

	return client.service.Logout(ctx, req)
}

func (client *AuthClient) SetMasterKey(accessToken string, salt, verifier []byte) (*pb.SetMasterKeyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.SetMasterKeyRequest{
		Salt:     salt,
		Verifier: verifier,
	}

	return client.service.SetMasterKey(ctx, req)
}

func (client *AuthClient) GetMasterKeyData(accessToken string) (*pb.GetMasterKeyDataResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.GetMasterKeyDataRequest{}

	return client.service.GetMasterKeyData(ctx, req)
}

func (client *AuthClient) HasMasterKey(accessToken string) (*pb.HasMasterKeyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.HasMasterKeyRequest{}

	return client.service.HasMasterKey(ctx, req)
}
