package response

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/results"
)

type ProductResponseMapper interface {
	ToCreateProductReponse(_ context.Context, resp interface{}) (interface{}, error)
	ToUpdateProductResponse(_ context.Context, resp interface{}) (interface{}, error)
	ToAddProductCategoriesResponse(_ context.Context, resp interface{}) (interface{}, error)
	ToAddProductVariationsResponse(_ context.Context, resp interface{}) (interface{}, error)
}

type ProductResponseMapperImpl struct{}

func (p *ProductResponseMapperImpl) ToCreateProductReponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*results.CreateProductResult)
	variations := make([]string, len(res.Variations))
	for idx, variation := range res.Variations {
		variations[idx] = variation.Id.Name
	}
	categoryIds := make([]string, len(res.CategoryIds))
	for idx, categoryId := range res.CategoryIds {
		categoryIds[idx] = string(categoryId)
	}
	imageRes := make([]*proto.ImageData, len(res.Images))
	for idx, image := range res.Images {
		imageRes[idx] = &proto.ImageData{
			MetaData: &proto.ImageArg{
				Size:  int32(image.Size),
				Width: int32(image.Width),
				Heigh: int32(image.Height),
				Type:  image.Type,
			},
			Url:    image.Id.Url.Key,
			Method: "PUT",
		}
	}
	return &proto.ProductData{
		Id:          string(res.Id),
		Name:        res.Name,
		Description: res.Description,
		Unit:        res.Unit,
		ShopId:      string(res.ShopId),
		Variations:  variations,
		CategoryIds: categoryIds,
		Images:      imageRes,
	}, nil
}

func (p *ProductResponseMapperImpl) ToUpdateProductResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*results.UpdateProductResult)
	variations := make([]string, len(res.Variations))
	for idx, variation := range res.Variations {
		variations[idx] = variation.Id.Name
	}
	categoryIds := make([]string, len(res.CategoryIds))
	for idx, categoryId := range res.CategoryIds {
		categoryIds[idx] = string(categoryId)
	}
	return &proto.ProductData{
		Id:          string(res.Id),
		Name:        res.Name,
		Description: res.Description,
		Unit:        res.Unit,
		ShopId:      string(res.ShopId),
		Variations:  variations,
		CategoryIds: categoryIds,
	}, nil
}
func (p *ProductResponseMapperImpl) ToAddProductCategoriesResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*results.AddProductCategoriesResult)
	variations := make([]string, len(res.Variations))
	for idx, variation := range res.Variations {
		variations[idx] = variation.Id.Name
	}
	categoryIds := make([]string, len(res.CategoryIds))
	for idx, categoryId := range res.CategoryIds {
		categoryIds[idx] = string(categoryId)
	}
	return &proto.ProductData{
		Id:          string(res.Id),
		Name:        res.Name,
		Description: res.Description,
		Unit:        res.Unit,
		ShopId:      string(res.ShopId),
		Variations:  variations,
		CategoryIds: categoryIds,
	}, nil
}
func (p *ProductResponseMapperImpl) ToAddProductVariationsResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*results.AddProductVariationsResult)
	variations := make([]string, len(res.Variations))
	for idx, variation := range res.Variations {
		variations[idx] = variation.Id.Name
	}
	categoryIds := make([]string, len(res.CategoryIds))
	for idx, categoryId := range res.CategoryIds {
		categoryIds[idx] = string(categoryId)
	}
	return &proto.ProductData{
		Id:          string(res.Id),
		Name:        res.Name,
		Description: res.Description,
		Unit:        res.Unit,
		ShopId:      string(res.ShopId),
		Variations:  variations,
		CategoryIds: categoryIds,
	}, nil
}

func NewProductResponseMapper() ProductResponseMapper {
	return &ProductResponseMapperImpl{}
}
