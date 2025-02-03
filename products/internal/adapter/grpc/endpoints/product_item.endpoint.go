package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/grpc/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
)

type ProductItemEndpoints struct {
	CreateProductItem endpoint.Endpoint
}

func MakeProductItemEndpoints(c handlers.ProductItemCommandHandler, tracer opentracing.Tracer) ProductItemEndpoints {
	return ProductItemEndpoints{
		CreateProductItem: middlewares.ExtractTracingMiddleware(tracer, "Create product item")(makeCreateProductItemEndpoint(c)),
	}
}

func makeCreateProductItemEndpoint(c handlers.ProductItemCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.CreateProductItemCommand)
		res, err := c.CreateProductItem(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
