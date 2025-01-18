package valueobjects

type CategoryId string

func NewCategoryId(id string) CategoryId {
	return CategoryId(id)
}

func (c CategoryId) IsEqual(o CategoryId) bool {
	return string(c) == string(o)
}
