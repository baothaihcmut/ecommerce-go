package commands

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
)

type Address struct {
	Priority int
	Street   string
	Town     string
	City     string
	Province string
}
type CustomerInfo struct {
}

type ShopOwnerInfo struct {
	BussinessLincese string
}
type SignUpCommand struct {
	Email         string
	Password      string
	PhoneNumber   string
	Addresses     []user.AddressArg
	FirstName     string
	LastName      string
	Role          enums.Role
	CustomerInfo  *CustomerInfo
	ShopOwnerInfo *ShopOwnerInfo
}
