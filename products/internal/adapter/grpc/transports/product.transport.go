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
	createProductHandler        grpc.Handler
	updateProductHandler        grpc.Handler
	addProductCategoriesHandler grpc.Handler
	addProductVariationsHandler grpc.Handler
}

// AddProductCategories implements proto.ProductServiceServer.
func (p *ProductServer) AddProductCategories(ctx context.Context, req *proto.AddProductCategoriesRequest) (*proto.UpdateProductResponse, error) {
	_, res, err := p.addProductCategoriesHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.UpdateProductResponse{
			Data: nil,
			Status: &proto.Status{
				Message: err.Error(),
				Code:    errors.MapGrpcErrorCode(err).String(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.UpdateProductResponse{
		Data: res.(*proto.ProductData),
		Status: &proto.Status{
			Message: "create product success",
			Code:    codes.OK.String(),
		},
	}, nil
}

// AddProductVariations implements proto.ProductServiceServer.
func (p *ProductServer) AddProductVariations(ctx context.Context, req *proto.AddProductVariationsRequest) (*proto.AddProductVariationsResponse, error) {
	_, resp, err := p.addProductVariationsHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.AddProductVariationsResponse{
			Data: nil,
			Status: &proto.Status{
				Message: err.Error(),
				Code:    errors.MapGrpcErrorCode(err).String(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.AddProductVariationsResponse{
		Data: resp.(*proto.ProductData),
		Status: &proto.Status{
			Message: "create product success",
			Code:    codes.OK.String(),
		},
	}, nil
}

// UpdateProduct implements proto.ProductServiceServer.
func (p *ProductServer) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	_, resp, err := p.updateProductHandler.ServeGRPC(ctx, req)
	if err != nil {
		return &proto.UpdateProductResponse{
			Data: nil,
			Status: &proto.Status{
				Message: err.Error(),
				Code:    errors.MapGrpcErrorCode(err).String(),
				Details: []string{err.Error()},
			},
		}, nil
	}
	return &proto.UpdateProductResponse{
		Data: resp.(*proto.ProductData),
		Status: &proto.Status{
			Message: "create product success",
			Code:    codes.OK.String(),
		},
	}, nil
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
		Data: resp.(*proto.ProductData),
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
		updateProductHandler: grpc.NewServer(
			endpoints.UpdateProduct,
			requestMapper.ToUpdateProductCommand,
			responseMapper.ToUpdateProductResponse,
		),
		addProductCategoriesHandler: grpc.NewServer(
			endpoints.AddProductCategories,
			requestMapper.ToAddProductCategoriesCommand,
			responseMapper.ToAddProductCategoriesResponse,
		),
		addProductVariationsHandler: grpc.NewServer(
			endpoints.AddProductVariations,
			requestMapper.ToAddProductVariationsCommand,
			responseMapper.ToAddProductVariationsResponse,
		),
	}
}
