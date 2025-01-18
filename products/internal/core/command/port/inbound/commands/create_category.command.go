package commands

type CreateCategoryCommand struct {
	Name              string
	ParentCategoryIds []string
}
