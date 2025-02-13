package response

type BulkCreateCategoriesResponseDTO struct {
	Categories []*CreateCategoryResponseDTO `json:"categories"`
}
