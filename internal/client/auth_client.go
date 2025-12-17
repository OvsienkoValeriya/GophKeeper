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

func (client *AuthClient) Login(username, password string) (userId, accessToken, refreshToken string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: proto.String(username),
		Password: proto.String(password),
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", "", "", err
	}

	return *res.UserId, *res.AccessToken, *res.RefreshToken, nil
}

func (client *AuthClient) Register(username, password string) (userId, accessToken, refreshToken string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.RegisterRequest{
		Username: proto.String(username),
		Password: proto.String(password),
	}

	res, err := client.service.Register(ctx, req)
	if err != nil {
		return "", "", "", err
	}

	return *res.UserId, *res.AccessToken, *res.RefreshToken, nil
}

func (client *AuthClient) RefreshToken(accessToken string) (newAccessToken, newRefreshToken string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.RefreshTokenRequest{}

	res, err := client.service.RefreshToken(ctx, req)
	if err != nil {
		return "", "", err
	}

	return *res.AccessToken, *res.RefreshToken, nil
}

func (client *AuthClient) Logout(accessToken string) (res bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.LogoutRequest{
		AccessToken: proto.String(accessToken),
	}

	resp, err := client.service.Logout(ctx, req)
	if err != nil {
		return false, err
	}

	return *resp.Success, nil
}

func (client *AuthClient) SetMasterKey(accessToken string, salt, verifier []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.SetMasterKeyRequest{
		Salt:     salt,
		Verifier: verifier,
	}

	_, err := client.service.SetMasterKey(ctx, req)
	return err
}

func (client *AuthClient) GetMasterKeyData(accessToken string) (salt, verifier []byte, hasMasterKey bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.GetMasterKeyDataRequest{}

	res, err := client.service.GetMasterKeyData(ctx, req)
	if err != nil {
		return nil, nil, false, err
	}

	return res.Salt, res.Verifier, res.GetHasMasterKey(), nil
}

func (client *AuthClient) HasMasterKey(accessToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	req := &pb.HasMasterKeyRequest{}

	res, err := client.service.HasMasterKey(ctx, req)
	if err != nil {
		return false, err
	}

	return res.GetHasMasterKey(), nil
}
