package entities

import (
	"errors"

	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/customer/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
)

var (
	EmailRequiredErr = errors.New("email is required")
	NameRequiredErr  = errors.New("last name is required")
	MinAddressErr    = errors.New("at least one address is required")
)

type User struct {
	Id          valueobject.UserId
	Email       valueobject.Email
	PhoneNumber valueobject.PhoneNumber
	Address     []valueobject.Address
	Role        enums.Role
	FirstName   string
	LastName    string
}

func NewUser(
	id valueobject.UserId,
	email valueobject.Email,
	phoneNumber valueobject.PhoneNumber,
	address []valueobject.Address,
	role enums.Role,
	firstName string,
	lastName string,
) (*User, error) {

	if email == "" {
		return nil, EmailRequiredErr
	}
	if firstName == "" {
		return nil, NameRequiredErr
	}
	if lastName == "" {
		return nil, NameRequiredErr
	}
	if len(address) == 0 {
		return nil, MinAddressErr
	}

	// Create and return a new User instance
	return &User{
		Id:          id,
		Email:       email,
		PhoneNumber: phoneNumber,
		Address:     address,
		Role:        role,
		FirstName:   firstName,
		LastName:    lastName,
	}, nil
}
