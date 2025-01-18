package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/errors"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/mappers/request"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/mappers/response"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/proto"
	gt "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
)

type CategoryServer struct {
	CreateCategoryHandler  gt.Handler
	FinaAllCategoryHandler gt.Handler
}

// CreateCategory implements proto.CategoryServiceServer.
func (c *CategoryServer) CreateCategory(ctx context.Context, req *proto.CreateCategoryRequest) (*proto.CreateCategoryResponse, error) {
	_, res, err := c.CreateCategoryHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.CreateCategoryResponse{
			Data: nil,
			Status: &proto.Status{
				Code:    errors.MapGrpcErrorCode(err).String(),
				Message: err.Error(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.CreateCategoryResponse{
		Data: res.(*proto.CategoryData),
		Status: &proto.Status{
			Code:    codes.OK.String(),
			Message: "Create category sucess",
		},
	}, nil
}

// FindAllCategory implements proto.CategoryServiceServer.
func (c *CategoryServer) FindAllCategory(ctx context.Context, req *proto.FindAllCategoryRequest) (*proto.FindAllCategoryResponse, error) {
	_, res, err := c.FinaAllCategoryHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.FindAllCategoryResponse{
			Data: nil,
			Status: &proto.Status{
				Code:    errors.MapGrpcErrorCode(err).String(),
				Message: err.Error(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.FindAllCategoryResponse{
		Data: res.(*proto.FindAllCategoryData),
		Status: &proto.Status{
			Code:    codes.OK.String(),
			Message: "Find all category sucess",
		},
	}, nil
}

func NewCategoryServer(
	endpoints endpoints.CategoryEndpoints,
	requestMapper request.CategoryRequestMapper,
	responseMapper response.CategoryResponseMapper,
) proto.CategoryServiceServer {
	return &CategoryServer{
		CreateCategoryHandler: gt.NewServer(
			endpoints.CreateCategory,
			requestMapper.ToCreateCategoryCommand,
			responseMapper.ToCreateCategoryResponse,
		),
		FinaAllCategoryHandler: gt.NewServer(
			endpoints.FindAllCategory,
			requestMapper.ToFindAllCategoryQuery,
			responseMapper.ToFindAllCategoryResponse,
		),
	}
}
