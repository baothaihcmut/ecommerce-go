package request

type CreateCategoryRequestDTO struct {
	Name              string   `json:"name" validate:"required"`
	ParentCategoryIds []string `json:"parent_category_ids" validate:"required"`
}
