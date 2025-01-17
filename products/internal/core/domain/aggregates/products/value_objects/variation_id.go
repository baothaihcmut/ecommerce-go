package valueobjects

type VariationId struct {
	ProductId ProductId
	Name      string
}

func NewVariationId(productId ProductId, name string) VariationId {
	return VariationId{
		ProductId: productId,
		Name:      name,
	}
}
func (v VariationId) IsEqual(o VariationId) bool {
	return v.ProductId.IsEqual(o.ProductId) && v.Name == o.Name
}
