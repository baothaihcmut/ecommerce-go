package valueobjects

type ProductId string

func NewProductId(id string) ProductId {
	return ProductId(id)
}
func (p ProductId) IsEqual(o ProductId) bool {
	return string(p) == string(o)
}
