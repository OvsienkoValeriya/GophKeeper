package services

import (
	"context"
	"errors"

	pb "github.com/OvsienkoValeriya/GophKeeper/api/gen"
	"github.com/OvsienkoValeriya/GophKeeper/internal/models"
	"github.com/OvsienkoValeriya/GophKeeper/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type ResourceServer struct {
	pb.UnimplementedResourceServiceServer
	resourceService *service.ResourceService
}

func NewResourceServer(resourceService *service.ResourceService) *ResourceServer {
	return &ResourceServer{
		resourceService: resourceService,
	}
}

func (s *ResourceServer) CreateResource(ctx context.Context, req *pb.CreateResourceRequest) (*pb.CreateResourceResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	resourceType := models.ResourceType(req.GetType())
	if !isValidResourceType(resourceType) {
		return nil, status.Error(codes.InvalidArgument, "invalid resource type")
	}

	resource, err := s.resourceService.Upload(ctx, userID, req.GetName(), resourceType, req.GetData())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create resource: %v", err)
	}

	return &pb.CreateResourceResponse{
		Id:        proto.Int64(resource.ID),
		Name:      proto.String(resource.Name),
		Type:      proto.String(string(resource.Type)),
		Size:      proto.Int64(resource.Size),
		CreatedAt: proto.String(resource.CreatedAt.Format("2006-01-02T15:04:05Z")),
	}, nil
}

func (s *ResourceServer) GetResource(ctx context.Context, req *pb.GetResourceRequest) (*pb.GetResourceResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	resource, data, err := s.resourceService.Get(ctx, userID, req.GetId())
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		if errors.Is(err, service.ErrResourceNotFound) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get resource: %v", err)
	}

	return &pb.GetResourceResponse{
		Id:        proto.Int64(resource.ID),
		Name:      proto.String(resource.Name),
		Type:      proto.String(string(resource.Type)),
		Data:      data,
		Size:      proto.Int64(resource.Size),
		CreatedAt: proto.String(resource.CreatedAt.Format("2006-01-02T15:04:05Z")),
		UpdatedAt: proto.String(resource.UpdatedAt.Format("2006-01-02T15:04:05Z")),
	}, nil
}

func (s *ResourceServer) GetResourceByName(ctx context.Context, req *pb.GetResourceByNameRequest) (*pb.GetResourceResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	resource, data, err := s.resourceService.GetByName(ctx, userID, req.GetName())
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		if errors.Is(err, service.ErrResourceNotFound) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get resource: %v", err)
	}

	return &pb.GetResourceResponse{
		Id:        proto.Int64(resource.ID),
		Name:      proto.String(resource.Name),
		Type:      proto.String(string(resource.Type)),
		Data:      data,
		Size:      proto.Int64(resource.Size),
		CreatedAt: proto.String(resource.CreatedAt.Format("2006-01-02T15:04:05Z")),
		UpdatedAt: proto.String(resource.UpdatedAt.Format("2006-01-02T15:04:05Z")),
	}, nil
}

func (s *ResourceServer) ListResources(ctx context.Context, req *pb.ListResourcesRequest) (*pb.ListResourcesResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	resources, err := s.resourceService.GetAll(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list resources: %v", err)
	}

	pbResources := make([]*pb.GetResourceResponse, len(resources))
	for i, r := range resources {
		pbResources[i] = &pb.GetResourceResponse{
			Id:        proto.Int64(r.ID),
			Name:      proto.String(r.Name),
			Type:      proto.String(string(r.Type)),
			Size:      proto.Int64(r.Size),
			CreatedAt: proto.String(r.CreatedAt.Format("2006-01-02T15:04:05Z")),
			UpdatedAt: proto.String(r.UpdatedAt.Format("2006-01-02T15:04:05Z")),
		}
	}

	return &pb.ListResourcesResponse{
		Resources: pbResources,
	}, nil
}

func (s *ResourceServer) DeleteResource(ctx context.Context, req *pb.DeleteResourceRequest) (*pb.DeleteResourceResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	err = s.resourceService.Delete(ctx, userID, req.GetId())
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete resource: %v", err)
	}

	return &pb.DeleteResourceResponse{
		Success: proto.Bool(true),
	}, nil
}

func getUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0, errors.New("userID not found in context")
	}
	return userID, nil
}

func isValidResourceType(t models.ResourceType) bool {
	switch t {
	case models.TypeCredentials, models.TypeText, models.TypeBinary, models.TypeCard:
		return true
	default:
		return false
	}
}
