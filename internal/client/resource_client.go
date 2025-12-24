package client

import (
	"context"
	"time"

	pb "github.com/OvsienkoValeriya/GophKeeper/api/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type ResourceClient struct {
	service    pb.ResourceServiceClient
	tokenStore *FileTokenStore
}

func NewResourceClient(cc *grpc.ClientConn, tokenStore *FileTokenStore) *ResourceClient {
	service := pb.NewResourceServiceClient(cc)
	return &ResourceClient{
		service:    service,
		tokenStore: tokenStore,
	}
}

func (c *ResourceClient) withAuth(ctx context.Context) context.Context {
	accessToken, _, err := c.tokenStore.LoadTokens()
	if err != nil || accessToken == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)
}

// CreateResource creates a new resource
// Parameters:
//   - name: name of the resource
//   - resourceType: type of the resource
//   - encryptedData: encrypted data of the resource
//
// Returns:
//   - int64: id of the created resource
//   - error: error if the resource creation failed
func (c *ResourceClient) CreateResource(name, resourceType string, encryptedData []byte) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = c.withAuth(ctx)

	req := &pb.CreateResourceRequest{
		Name: proto.String(name),
		Type: proto.String(resourceType),
		Data: encryptedData,
	}

	res, err := c.service.CreateResource(ctx, req)
	if err != nil {
		return 0, err
	}

	return res.GetId(), nil
}

// GetResourceByName gets a resource by name
// Parameters:
//   - name: name of the resource
//
// Returns:
//   - *pb.GetResourceResponse: resource information
//   - error: error if the resource retrieval failed
func (c *ResourceClient) GetResourceByName(name string) (*pb.GetResourceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	ctx = c.withAuth(ctx)

	req := &pb.GetResourceByNameRequest{
		Name: proto.String(name),
	}

	return c.service.GetResourceByName(ctx, req)
}

// GetResource gets a resource by id
// Parameters:
//   - id: id of the resource
//
// Returns:
//   - *pb.GetResourceResponse: resource information
//   - error: error if the resource retrieval failed
func (c *ResourceClient) GetResource(id int64) (*pb.GetResourceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = c.withAuth(ctx)

	req := &pb.GetResourceRequest{
		Id: proto.Int64(id),
	}

	return c.service.GetResource(ctx, req)
}

// ListResources lists all resources
// Parameters:
//   - type: type of the resources (credentials, text, binary, card)
//
// Returns:
//   - *pb.ListResourcesResponse: list of resources
//   - error: error if the resource listing failed
func (c *ResourceClient) ListResources() (*pb.ListResourcesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = c.withAuth(ctx)

	req := &pb.ListResourcesRequest{}

	return c.service.ListResources(ctx, req)
}

// DeleteResource deletes a resource by id
// Parameters:
//   - id: id of the resource
//
// Returns:
//   - error: error if the resource deletion failed
func (c *ResourceClient) DeleteResource(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = c.withAuth(ctx)

	req := &pb.DeleteResourceRequest{
		Id: proto.Int64(id),
	}

	_, err := c.service.DeleteResource(ctx, req)
	return err
}
