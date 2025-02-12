package response

type CreateCategoryResponseDTO struct {
	Id                string   `json:"id"`
	Name              string   `json:"name"`
	ParentCategoryIds []string `json:"parent_category_ids"`
}
