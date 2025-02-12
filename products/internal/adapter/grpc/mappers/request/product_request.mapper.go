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
	ToUpdateProductCommand(_ context.Context, request interface{}) (interface{}, error)
	ToAddProductCategoriesCommand(_ context.Context, request interface{}) (interface{}, error)
	ToAddProductVariationsCommand(_ context.Context, request interface{}) (interface{}, error)
}

type ProductRequestMapperImpl struct {
}

func (p *ProductRequestMapperImpl) ToCreateProductCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.CreateProductRequest)
	categoryIds := make([]categoryValueobjects.CategoryId, len(req.CategoryIds))
	for idx, categoryId := range req.CategoryIds {
		categoryIds[idx] = categoryValueobjects.NewCategoryId(categoryId)
	}
	imageArgs := make([]commands.CreateImageCommand, len(req.Images))
	for idx, image := range req.Images {
		imageArgs[idx] = commands.CreateImageCommand{
			Size:   int(image.Size),
			Width:  int(image.Width),
			Height: int(image.Heigh),
			Type:   image.Type,
		}
	}
	return &commands.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		ShopId:      valueobjects.NewShopId(req.ShopId),
		CategoryIds: categoryIds,
		Variations:  req.Variations,
		Images:      imageArgs,
	}, nil

}
func (p *ProductRequestMapperImpl) ToUpdateProductCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.UpdateProductRequest)
	cmd := commands.UpdateProductCommand{}
	if req.Name != nil {
		cmd.Name = &req.Name.Value
	}
	if req.Description != nil {
		cmd.Description = &req.Description.Value
	}
	if req.Unit != nil {
		cmd.Unit = &req.Unit.Value
	}
	cmd.Id = valueobjects.NewProductId(req.Id)

	return &cmd, nil
}

func (p *ProductRequestMapperImpl) ToAddProductCategoriesCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.AddProductCategoriesRequest)
	newCategories := make([]categoryValueobjects.CategoryId, len(req.NewCategoryIds))
	for idx, categoryId := range req.NewCategoryIds {
		newCategories[idx] = categoryValueobjects.CategoryId(categoryId)
	}
	return &commands.AddProductCategoiesCommand{
		ProductId:     valueobjects.ProductId(req.ProductId),
		NewCategories: newCategories,
	}, nil
}
func (p *ProductRequestMapperImpl) ToAddProductVariationsCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.AddProductVariationsRequest)
	productId := valueobjects.NewProductId(req.ProductId)

	return &commands.AddProductVariationsCommand{
		ProductId:     productId,
		NewVariations: req.NewVariations,
	}, nil
}

func NewProductRequestMapper() ProductRequestMapper {
	return &ProductRequestMapperImpl{}
}
