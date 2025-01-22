package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/go-kit/kit/endpoint"
)

type ProductEndpoints struct {
	CreateProduct endpoint.Endpoint
}

func MakeProductEndpoints(c commandHandler.ProductCommandHandler) ProductEndpoints {
	return ProductEndpoints{
		CreateProduct: makeCreateProductEndpoint(c),
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
