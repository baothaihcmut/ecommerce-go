package user

import (
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/google/uuid"
)

var (
	EmailRequiredErr = errors.New("email is required")
	NameRequiredErr  = errors.New("last name is required")
	MinAddressErr    = errors.New("at least one address is required")
)

type User struct {
	Id          valueobject.UserId
	Email       valueobject.Email
	Password    valueobject.Password
	PhoneNumber valueobject.PhoneNumber
	Address     []valueobject.Address
	Role        enums.Role
	FirstName   string
	LastName    string
	Customer    *entities.Customer
	ShopOwner   *entities.ShopOwner
}

func validate(
	email valueobject.Email,
	address []valueobject.Address,
	firstName string,
	lastName string,
) error {
	if email == "" {
		return EmailRequiredErr
	}
	if firstName == "" {
		return NameRequiredErr
	}
	if lastName == "" {
		return NameRequiredErr
	}
	if len(address) == 0 {
		return MinAddressErr
	}
	return nil
}

func NewCustomer(
	email valueobject.Email,
	password valueobject.Password,
	phoneNumber valueobject.PhoneNumber,
	address []valueobject.Address,
	firstName string,
	lastName string,
) (*User, error) {

	err := validate(email, address, firstName, lastName)
	id, err := valueobject.NewUserId(uuid.New())
	if err != nil {
		return nil, err
	}
	customer, err := entities.NewCustomer()
	if err != nil {
		return nil, err
	}

	// Create and return a new User instance
	return &User{
		Id:          *id,
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
		Address:     address,
		Role:        enums.CUSTOMER,
		FirstName:   firstName,
		LastName:    lastName,
		Customer:    customer,
	}, nil
}

func NewShopOwner(
	email valueobject.Email,
	passwod valueobject.Password,
	phoneNumber valueobject.PhoneNumber,
	address []valueobject.Address,
	firstName string,
	lastName string,
	bussinessLincese string,
) (*User, error) {
	err := validate(email, address, firstName, lastName)
	id, err := valueobject.NewUserId(uuid.New())
	if err != nil {
		return nil, err
	}
	shopOwner := entities.NewShopOwner(bussinessLincese)
	return &User{
		Id:          *id,
		Email:       email,
		Password:    passwod,
		PhoneNumber: phoneNumber,
		Address:     address,
		Role:        enums.SHOP_OWNER,
		FirstName:   firstName,
		LastName:    lastName,
		ShopOwner:   shopOwner,
	}, nil

}
