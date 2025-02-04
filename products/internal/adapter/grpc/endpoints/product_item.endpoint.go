package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/trace"
)

type ProductItemEndpoints struct {
	CreateProductItem endpoint.Endpoint
}

func MakeProductItemEndpoints(c handlers.ProductItemCommandHandler, tracer trace.Tracer) ProductItemEndpoints {
	return ProductItemEndpoints{
		CreateProductItem: makeCreateProductItemEndpoint(c, tracer),
	}
}

func makeCreateProductItemEndpoint(c handlers.ProductItemCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracer.Start(ctx, "ProductItem.Create: endpoint")
		defer span.End()
		req := request.(*commands.CreateProductItemCommand)
		res, err := c.CreateProductItem(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
