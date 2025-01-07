package dtos

type CustomerInfo struct{}

type ShopOwnerInfo struct {
	BusinessLicense string `json:"business_name" validate:"required"`
}

type CreateUserDto struct {
	Email         string         `json:"email" validate:"required,email"`
	Password      string         `json:"password" validate:"required"`
	FirstName     string         `json:"first_name" validate:"required"`
	LastName      string         `json:"last_name" validate:"required"`
	PhoneNumber   string         `json:"phone" validate:"required"`
	Role          Role           `json:"role" validate:"required,oneof=USER SHOP_OWNER ADMIN"`
	Addresses     []Address      `json:"addresses" validate:"required"`
	ShopOwnerInfo *ShopOwnerInfo `json:"shop_owner_info"`
	CustomerInfo  *CustomerInfo  `json:"customer_info"`
}
