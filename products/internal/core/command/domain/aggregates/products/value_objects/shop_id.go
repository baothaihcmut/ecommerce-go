package valueobjects

type ShopId string

func NewShopId(Id string) ShopId {
	return ShopId(Id)
}
