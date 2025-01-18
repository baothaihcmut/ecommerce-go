package response

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/proto"
	commandResult "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
	queryResult "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/results"
)

type CategoryResponseMapper interface {
	ToCreateCategoryResponse(context.Context, interface{}) (interface{}, error)
	ToFindAllCategoryResponse(context.Context, interface{}) (interface{}, error)
}

type CategoryResponseMapperImpl struct{}

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
