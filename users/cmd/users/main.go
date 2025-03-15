package main

import (
	"net"
	"os"
	"time"

	userProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	grpcListener, listErr := net.Listen("tcp", ":50051")
	if listErr != nil {
		os.Exit(1)
	}
	// grpc options
	serverOptions := []grpc.ServerOption{
		// Unary option
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Duration(5),
			MaxConnectionAge:  10 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           20 * time.Second,
		}),
	}
	baseServer := grpc.NewServer(serverOptions...)
	userProto.RegisterAuthServiceServer(baseServer, services.NewAuthService())
	baseServer.Serve(grpcListener)
}
