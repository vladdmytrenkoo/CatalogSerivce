package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"CatalogService/internal/pkg/config"
	"CatalogService/internal/pkg/migration"
	"CatalogService/internal/services"
	"CatalogService/internal/transport/grpc/product"
)

func main() {
	cfg, err := config.InitConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := migration.RunMigrations(cfg.Database.SpannerDB, cfg.App.MigrationDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// If run with --migrate-only flag, exit after migrations
	for _, arg := range os.Args {
		if arg == "--migrate-only" {
			log.Println("Migration-only mode: exiting")
			return
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := spanner.NewClient(ctx, cfg.Database.SpannerDB,
		option.WithEndpoint("localhost:9010"),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		log.Fatalf("failed to create spanner client: %v", err)
	}
	defer client.Close()

	opts := services.NewOptions(client)
	handler := product.NewHandler(opts.Commands, opts.Queries)

	grpcServer := grpc.NewServer()
	handler.Register(grpcServer)
	reflection.Register(grpcServer)

	product.ListGRPCMethods(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	errChan := make(chan error, 1)
	go func() {
		log.Printf("Starting gRPC server on port %d", cfg.App.Port)
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	case err := <-errChan:
		log.Fatalf("gRPC server error: %v", err)
	}
}
