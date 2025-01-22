package response

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
	commandResult "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
	queryResult "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/results"
)

type CategoryResponseMapper interface {
	ToCreateCategoryResponse(context.Context, interface{}) (interface{}, error)
	ToBulkCreateCategoriesResponse(context.Context, interface{}) (interface{}, error)
	ToFindAllCategoryResponse(context.Context, interface{}) (interface{}, error)
}

type CategoryResponseMapperImpl struct{}

// ToBulkCreateCategoriesResponse implements CategoryResponseMapper.
func (c *CategoryResponseMapperImpl) ToBulkCreateCategoriesResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*results.BulkCreateCategoriesResult)
	categoriesResponse := make([]*proto.CategoryData, len(res.Categories))
	for idx, category := range res.Categories {
		parentCategoryIds := make([]string, len(category.ParentCategoryId))
		for idx, val := range category.ParentCategoryId {
			parentCategoryIds[idx] = string(val)
		}
		categoriesResponse[idx] = &proto.CategoryData{
			Id:                string(category.Id),
			Name:              category.Name,
			ParentCategoryIds: parentCategoryIds,
		}
	}
	return categoriesResponse, nil
}

func (c *CategoryResponseMapperImpl) ToCreateCategoryResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*commandResult.CreateCategoryResult)
	parentCategoryIds := make([]string, len(res.ParentCategoryId))
	for idx, val := range res.ParentCategoryId {
		parentCategoryIds[idx] = string(val)
	}
	return &proto.CategoryData{
		Id:                string(res.Id),
		Name:              res.Name,
		ParentCategoryIds: parentCategoryIds,
	}, nil
}

// ToFindAllCategoryResponse implements CategoryResponseMapper.
func (c *CategoryResponseMapperImpl) ToFindAllCategoryResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*queryResult.FindAllCategoryResult)
	items := make([]*proto.CategoryData, len(res.Data))
	for idx, val := range res.Data {
		items[idx] = &proto.CategoryData{
			Id:                val.Id,
			Name:              val.Name,
			ParentCategoryIds: val.ParentCategoryIds,
		}
	}
	return &proto.FindAllCategoryData{
		Items: items,
		Pagination: &proto.PaginationMeta{
			CurrentPage:  int32(res.Pagination.CurrentPage),
			PageSize:     int32(res.Pagination.PageSize),
			TotalPage:    int32(res.Pagination.TotalPage),
			TotalElement: int32(res.Pagination.TotalItem),
		},
	}, nil
}

func NewCategoryResponseMapper() CategoryResponseMapper {
	return &CategoryResponseMapperImpl{}
}
