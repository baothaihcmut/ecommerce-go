package queries

import (
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/filter"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/sort"
)

type FindAllCategoryQuery struct {
	Filter     []filter.FilterParam
	Sort       []sort.SortParam
	Pagination pagination.PaginationParam
}
