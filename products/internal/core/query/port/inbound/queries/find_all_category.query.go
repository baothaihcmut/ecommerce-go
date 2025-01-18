package queries

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/filter"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/sort"
)

type FindAllCategoryQuery struct {
	Filter     []filter.FilterParam
	Sort       []sort.SortParam
	Pagination pagination.PaginationParam
}
