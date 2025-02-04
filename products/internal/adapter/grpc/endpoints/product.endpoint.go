package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/trace"
)

type ProductEndpoints struct {
	CreateProduct        endpoint.Endpoint
	UpdateProduct        endpoint.Endpoint
	AddProductCategories endpoint.Endpoint
	AddProductVariations endpoint.Endpoint
}

func MakeProductEndpoints(c commandHandler.ProductCommandHandler, tracer trace.Tracer) ProductEndpoints {
	return ProductEndpoints{
		CreateProduct:        makeCreateProductEndpoint(c, tracer),
		UpdateProduct:        makeUpdateProductEndpoint(c, tracer),
		AddProductCategories: makeAddProductCategoriesEndpoint(c, tracer),
		AddProductVariations: makeAddProductVariationsEndpoint(c, tracer),
	}
}

func makeCreateProductEndpoint(c commandHandler.ProductCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracer.Start(ctx, "Product.Create: endpoint")
		defer span.End()
		req := request.(*commands.CreateProductCommand)
		res, err := c.CreateProduct(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func makeUpdateProductEndpoint(c commandHandler.ProductCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracer.Start(ctx, "Product.Update: endpoint")
		defer span.End()
		req := request.(*commands.UpdateProductCommand)
		res, err := c.UpdateProduct(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeAddProductCategoriesEndpoint(c commandHandler.ProductCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracer.Start(ctx, "Product.AddCategories: endpoint")
		defer span.End()
		req := request.(*commands.AddProductCategoiesCommand)
		res, err := c.AddProductCategories(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeAddProductVariationsEndpoint(c commandHandler.ProductCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracer.Start(ctx, "Product.AddVariations: endpoints")
		defer span.End()
		req := request.(*commands.AddProductVariationsCommand)
		res, err := c.AddProductVariations(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
