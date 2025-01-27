package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/go-kit/kit/endpoint"
)

type ProductEndpoints struct {
	CreateProduct        endpoint.Endpoint
	UpdateProduct        endpoint.Endpoint
	AddProductCategories endpoint.Endpoint
	AddProductVariations endpoint.Endpoint
}

func MakeProductEndpoints(c commandHandler.ProductCommandHandler) ProductEndpoints {
	return ProductEndpoints{
		CreateProduct:        makeCreateProductEndpoint(c),
		UpdateProduct:        makeUpdateProductEndpoint(c),
		AddProductCategories: makeAddProductCategoriesEndpoint(c),
		AddProductVariations: makeAddProductVariationsEndpoint(c),
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
