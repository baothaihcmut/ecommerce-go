package request

type BulkCreateCategoriesRequestDTO struct {
	Categories []*CreateCategoryRequestDTO `json:"categories" validate:"required"`
}
