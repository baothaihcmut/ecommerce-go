package request

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/proto"
	categoryValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories/value_objects"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
)

type ProductRequestMapper interface {
	ToCreateProductCommand(_ context.Context, request interface{}) (interface{}, error)
}

type ProductRequestMapperImpl struct {
}

func (p ProductRequestMapperImpl) ToCreateProductCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.CreateProductRequest)
	categoryIds := make([]categoryValueobjects.CategoryId, len(req.CategoryIds))
	for idx, categoryId := range req.CategoryIds {
		categoryIds[idx] = categoryValueobjects.NewCategoryId(categoryId)
	}
	return &commands.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		ShopId:      valueobjects.NewShopId(req.ShopId),
		CategoryIds: categoryIds,
		Variations:  req.Variations,
	}, nil

}
func NewProductRequestMapper() ProductRequestMapper {
	return &ProductRequestMapperImpl{}
}
