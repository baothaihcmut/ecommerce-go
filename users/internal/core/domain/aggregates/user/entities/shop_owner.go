package entities

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
)

type ShopOwner struct {
	Id               valueobject.UserId
	BussinessLicense string
}

func NewShopOwner(id valueobject.UserId, bussinessLicense string) *ShopOwner {
	return &ShopOwner{
		Id:               id,
		BussinessLicense: bussinessLicense,
	}
}
