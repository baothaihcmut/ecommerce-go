package request

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/filter"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/sort"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/query/port/inbound/queries"
)

type CategoryRequestMapper interface {
	ToCreateCategoryCommand(_ context.Context, request interface{}) (interface{}, error)
	ToFindAllCategoryQuery(_ context.Context, request interface{}) (interface{}, error)
	ToBulkCreateCategoriesCommand(_ context.Context, request interface{}) (interface{}, error)
}

type CategoryRequestMapperImpl struct{}

func (m *CategoryRequestMapperImpl) ToCreateCategoryCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.CreateCategoryRequest)
	return &commands.CreateCategoryCommand{
		Name:              req.Name,
		ParentCategoryIds: req.ParentCategoryIds,
	}, nil
}

func (m *CategoryRequestMapperImpl) ToFindAllCategoryQuery(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.FindAllCategoryRequest)
	filters := make([]filter.FilterParam, len(req.FilterParams))
	for idx, val := range req.FilterParams {
		filters[idx] = filter.FilterParam{
			Field: val.Field,
			Value: val.Value,
		}
	}
	sorts := make([]sort.SortParam, len(req.SortParams))

	for idx, val := range req.SortParams {
		sorts[idx] = sort.SortParam{
			Field: val.Field,
		}
		if val.IsAsc {
			sorts[idx].Direction = sort.ASC
		} else {
			sorts[idx].Direction = sort.DESC
		}
	}
	return &queries.FindAllCategoryQuery{
		Filter: filters,
		Sort:   sorts,
		Pagination: pagination.PaginationParam{
			Page: int(req.PaginationParam.Page),
			Size: int(req.PaginationParam.Size),
		},
	}, nil
}

func (m *CategoryRequestMapperImpl) ToBulkCreateCategoriesCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.BulkCreateCategoryRequest)
	bulkCreateCommands := make([]*commands.CreateCategoryCommand, len(req.Categories))
	for idx, category := range req.Categories {
		bulkCreateCommands[idx] = &commands.CreateCategoryCommand{
			Name:              category.Name,
			ParentCategoryIds: category.ParentCategoryIds,
		}
	}
	return &commands.BulkCreateCategories{
		Categories: bulkCreateCommands,
	}, nil
}
func NewCategoryRequestMapper() CategoryRequestMapper {
	return &CategoryRequestMapperImpl{}
}
