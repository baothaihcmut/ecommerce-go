package entities

type ShopOwner struct {
	User             *User
	BussinessLicense string
}

func NewShopOwner(
	user *User,
	bussinessLicense string,
) *ShopOwner {
	return &ShopOwner{
		User:             user,
		BussinessLicense: bussinessLicense,
	}
}
