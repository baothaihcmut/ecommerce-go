package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/errors"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/mappers/request"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/mappers/response"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/proto"
	"github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
)

type ProductServer struct {
	createProductHandler grpc.Handler
}

// CreateProduct implements proto.ProductServiceServer.
func (p *ProductServer) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	_, resp, err := p.createProductHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.CreateProductResponse{
			Data: nil,
			Status: &proto.Status{
				Message: err.Error(),
				Code:    errors.MapGrpcErrorCode(err).String(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.CreateProductResponse{
		Data: resp.(*proto.CreateProductData),
		Status: &proto.Status{
			Message: "create product success",
			Code:    codes.OK.String(),
		},
	}, nil
}

func NewProductServer(
	endpoints endpoints.ProductEndpoints,
	requestMapper request.ProductRequestMapper,
	responseMapper response.ProductResponseMapper) proto.ProductServiceServer {
	return &ProductServer{
		createProductHandler: grpc.NewServer(
			endpoints.CreateProduct,
			requestMapper.ToCreateProductCommand,
			responseMapper.ToCreateProductReponse,
		),
	}
}
