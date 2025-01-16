package request

import "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/enums"

type Address struct {
	Priority int32  `json:"priority" validate:"gte=0"`          // Priority must be a non-negative integer
	Street   string `json:"street" validate:"required"`         // Street is required
	Town     string `json:"town" validate:"required"`           // Town is required
	City     string `json:"city" validate:"required"`           // City is required
	Province string `json:"province" validate:"required,alpha"` // Province must be alphabetic and required
}

type CustomerInfo struct {
}

type ShopOwnerInfo struct {
	BusinessLicense string `json:"business_license" validate:"required"`
}

type SignUpRequestDTO struct {
	Email         string         `json:"email" validate:"required,email"`
	Password      string         `json:"password" validate:"required,min=8"`
	PhoneNumber   string         `json:"phone_number" validate:"required"`
	Addresses     []Address      `json:"addresses" validate:"required,dive"`
	FirstName     string         `json:"first_name" validate:"required,alpha"`
	LastName      string         `json:"last_name" validate:"required,alpha"`
	Role          enums.Role     `json:"role" validate:"required,oneof=CUSTOMER SHOP_OWNER ADMIN"`
	CustomerInfo  *CustomerInfo  `json:"customer_info,omitempty" validate:"omitempty"`
	ShopOwnerInfo *ShopOwnerInfo `json:"shop_owner_info,omitempty" validate:"omitempty"`
}
