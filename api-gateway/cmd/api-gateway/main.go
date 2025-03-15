package main

import (
	"context"
	"net/http"

	userProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer conn.Close()
	mux := runtime.NewServeMux()
	err = userProto.RegisterAuthServiceHandler(context.Background(), mux, conn)
	if err != nil {
		return
	}
	http.ListenAndServe(":8080", mux)
}
