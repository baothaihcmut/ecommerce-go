package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/grpc/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
)

type ProductEndpoints struct {
	CreateProduct        endpoint.Endpoint
	UpdateProduct        endpoint.Endpoint
	AddProductCategories endpoint.Endpoint
	AddProductVariations endpoint.Endpoint
}

func MakeProductEndpoints(c commandHandler.ProductCommandHandler, tracer opentracing.Tracer) ProductEndpoints {
	return ProductEndpoints{
		CreateProduct:        middlewares.ExtractTracingMiddleware(tracer, "Create product")(makeCreateProductEndpoint(c)),
		UpdateProduct:        middlewares.ExtractTracingMiddleware(tracer, "Update product")(makeUpdateProductEndpoint(c)),
		AddProductCategories: middlewares.ExtractTracingMiddleware(tracer, "Add product categories")(makeAddProductCategoriesEndpoint(c)),
		AddProductVariations: middlewares.ExtractTracingMiddleware(tracer, "Add product variations")(makeAddProductVariationsEndpoint(c)),
	}
}

func makeCreateProductEndpoint(c commandHandler.ProductCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.CreateProductCommand)
		res, err := c.CreateProduct(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func makeUpdateProductEndpoint(c commandHandler.ProductCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.UpdateProductCommand)
		res, err := c.UpdateProduct(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeAddProductCategoriesEndpoint(c commandHandler.ProductCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.AddProductCategoiesCommand)
		res, err := c.AddProductCategories(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeAddProductVariationsEndpoint(c commandHandler.ProductCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.AddProductVariationsCommand)
		res, err := c.AddProductVariations(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
