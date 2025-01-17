package services

import (
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/grpc/tokens"
	"google.golang.org/grpc"
)

type GrpcConnectionService interface {
	GetConnection(tokens.GrpcToken, ...grpc.DialOption) (*grpc.ClientConn, error)
}

type GrpcConnectionServiceImpl struct {
	grpcAddr map[tokens.GrpcToken]string
}

func NewGrpcConnectionService(config config.GrpcServiceConfig) GrpcConnectionService {
	grpcAddr := map[tokens.GrpcToken]string{
		tokens.UserServiceToken: config.UserService,
	}
	return &GrpcConnectionServiceImpl{
		grpcAddr: grpcAddr,
	}
}

// GetConnection implements GrpcConnectionService.
func (g *GrpcConnectionServiceImpl) GetConnection(token tokens.GrpcToken, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.NewClient(g.grpcAddr[token], options...)
}
