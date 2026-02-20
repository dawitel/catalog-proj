package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/dawitel/catalog-proj/internal/pkg/config"
	"github.com/dawitel/catalog-proj/internal/pkg/logger"
	"github.com/dawitel/catalog-proj/internal/services"
	productv1 "github.com/dawitel/catalog-proj/proto/product/v1"
)

func main() {
	cfg := config.LoadFromEnv()

	ctx := context.Background()
	handler, client, err := services.NewProductHandler(ctx, services.Config{
		SpannerProject:  cfg.SpannerProject,
		SpannerInstance: cfg.SpannerInstance,
		SpannerDatabase: cfg.SpannerDatabase,
	})
	if err != nil {
		logger.Error("NewProductHandler failed", "err", err)
		os.Exit(1)
	}
	defer client.Close()

	srv := grpc.NewServer()
	productv1.RegisterProductServiceServer(srv, handler)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		logger.Error("listen failed", "err", err)
		os.Exit(1)
	}
	defer lis.Close()

	go func() {
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			logger.Error("serve failed", "err", err)
			os.Exit(1)
		}
	}()

	logger.Info("gRPC server listening", "addr", lis.Addr().String())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	logger.Info("shutdown signal received", "signal", sig.String())

	srv.GracefulStop()
	logger.Info("server stopped gracefully")
}
