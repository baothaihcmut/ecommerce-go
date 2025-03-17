package entities

type ShopOwner struct {
	BussinessLicense string
}

func NewShopOwner(
	bussinessLicense string,
) *ShopOwner {
	return &ShopOwner{
		BussinessLicense: bussinessLicense,
	}
}
