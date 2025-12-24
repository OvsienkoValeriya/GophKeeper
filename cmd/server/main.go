package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/OvsienkoValeriya/GophKeeper/api/gen"
	"github.com/OvsienkoValeriya/GophKeeper/internal/logger"
	"github.com/OvsienkoValeriya/GophKeeper/internal/repository/storage"
	"github.com/OvsienkoValeriya/GophKeeper/internal/server/auth"
	"github.com/OvsienkoValeriya/GophKeeper/internal/server/services"
	"github.com/OvsienkoValeriya/GophKeeper/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	databaseURL := getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/gophkeeper?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "your-super-secret-key-change-in-production")
	serverAddress := getEnv("SERVER_ADDRESS", ":50051")
	accessTokenDuration := 1 * time.Hour
	refreshTokenDuration := 7 * 24 * time.Hour

	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	minioBucket := getEnv("MINIO_BUCKET", "gophkeeper")
	minioUseSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	logger.InitDefault()
	defer logger.Sync()

	userStore, err := services.NewPostgresUserStore(databaseURL)
	if err != nil {
		logger.Sugar.Fatalf("Failed to connect to database: %v", err)
	}
	defer userStore.Close()
	logger.Sugar.Info("Connected to database")

	if err := userStore.RunMigrations(); err != nil {
		logger.Sugar.Fatalf("Failed to run migrations: %v", err)
	}
	logger.Sugar.Info("Migrations run successfully")

	resourceRepo, err := storage.NewPostgresResourceRepository(databaseURL)
	if err != nil {
		logger.Sugar.Fatalf("Failed to create resource repository: %v", err)
	}
	logger.Sugar.Info("Resource repository created")

	minioStorage, err := storage.NewMinioStorage(minioEndpoint, minioAccessKey, minioSecretKey, minioBucket, minioUseSSL)
	if err != nil {
		logger.Sugar.Fatalf("Failed to create MinIO storage: %v", err)
	}
	logger.Sugar.Info("MinIO storage connected")

	resourceService := service.NewResourceService(resourceRepo, minioStorage)

	jwtConfig := auth.NewJWTConfig(jwtSecret, accessTokenDuration, refreshTokenDuration)
	authServer := services.NewAuthServer(userStore, jwtConfig)
	resourceServer := services.NewResourceServer(resourceService)

	authInterceptor := services.NewAuthInterceptor(jwtConfig)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor()),
		grpc.StreamInterceptor(authInterceptor.StreamInterceptor()),
	)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterResourceServiceServer(grpcServer, resourceServer)
	reflection.Register(grpcServer)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup, grpcServer *grpc.Server) {
		defer wg.Done()
		<-ctx.Done()
		logger.Sugar.Info("Got signal to stop gRPC server")
		grpcServer.GracefulStop()
	}(ctx, &wg, grpcServer)

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		logger.Sugar.Fatalf("Failed to listen on %s: %v", serverAddress, err)
	}

	logger.Sugar.Infof("gRPC server listening on %s", serverAddress)
	if err := grpcServer.Serve(listener); err != nil {
		logger.Sugar.Fatalf("Failed to serve", "error", err)
	}
	wg.Wait()
	logger.Sugar.Info("gRPC server stopped")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
