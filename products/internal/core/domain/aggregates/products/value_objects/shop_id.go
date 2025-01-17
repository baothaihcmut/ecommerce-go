package valueobjects

type ShopId string

func NewShopId(id string) ShopId {
	return ShopId(id)
}
