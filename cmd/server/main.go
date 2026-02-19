package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/dawitel/product-catalog-service/internal/services"
	productv1 "github.com/dawitel/product-catalog-service/proto/product/v1"
)

func main() {
	cfg := services.Config{
		SpannerProject:  getEnv("SPANNER_PROJECT", "test-project"),
		SpannerInstance: getEnv("SPANNER_INSTANCE", "test-instance"),
		SpannerDatabase: getEnv("SPANNER_DATABASE", "product-catalog"),
	}
	port := getEnv("GRPC_PORT", "8080")
	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}

	ctx := context.Background()
	handler, client, err := services.NewProductHandler(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewProductHandler: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	srv := grpc.NewServer()
	productv1.RegisterProductServiceServer(srv, handler)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		os.Exit(1)
	}
	defer lis.Close()

	if err := srv.Serve(lis); err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		os.Exit(1)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
