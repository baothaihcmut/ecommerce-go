package valueobjects

type ProductItemId string

func NewProductItemId(id string) ProductItemId {
	return ProductItemId(id)
}
