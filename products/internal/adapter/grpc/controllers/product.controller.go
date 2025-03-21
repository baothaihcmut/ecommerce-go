package controllers

import (
	"context"
	"net/http"

	productProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/products/v1"
	v1 "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/shared/v1"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/grpc/mappers"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/inbound/handlers"
)

type ProductController struct {
	productProto.UnimplementedProductServiceServer
	productHandler handlers.ProductHandler
}

func NewProductController(productHandler handlers.ProductHandler) *ProductController {
	return &ProductController{
		productHandler: productHandler,
	}
}
func (p *ProductController) CreateProduct(ctx context.Context, req *productProto.CreateProductRequest) (*productProto.CreateProductResponse, error) {
	res, err := p.productHandler.CreateProduct(ctx, mappers.ToCreateProductCommand(req))
	if err != nil {
		return nil, err
	}
	return &productProto.CreateProductResponse{
		Status: &v1.Status{
			Message: "Create product success",
			Code:    http.StatusCreated,
		},
		Data: mappers.ToCreateProductReseponse(res),
	}, nil
}
