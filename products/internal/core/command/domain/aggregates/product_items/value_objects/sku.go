package valueobjects

type SKU string

func NewSKU(sku string) SKU {
	return SKU(sku)
}
