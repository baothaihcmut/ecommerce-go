package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		ctx, span := tracing.StartSpan(ctx, tracer, "ProductItem.Create: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*commands.CreateProductItemCommand)
		res, err := c.CreateProductItem(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
