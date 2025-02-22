package user

import (
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/google/uuid"
)

var (
	ErrEmailRequired        = errors.New("email is required")
	ErrNameRequired         = errors.New("last name is required")
	ErrMinAddress           = errors.New("at least one address is required")
	ErrInvalidTokenType     = errors.New("invalid token type")
	ErrMisMatchRefreshToken = errors.New("refresh token missmatch")
	ErrBadCredencial        = errors.New("email or password incorrect")
)

type AddressArg struct {
	Street   string
	Town     string
	City     string
	Province string
}
type ActivateShopOwnerArg struct {
	BussinessLincese string
}

type User struct {
	Id                  valueobject.UserId
	Email               valueobject.Email
	Password            valueobject.Password
	PhoneNumber         valueobject.PhoneNumber
	Address             []*entities.Address
	IsShopOwnerActive   bool
	FirstName           string
	LastName            string
	CurrentRefreshToken string
	Customer            *entities.Customer
	ShopOwner           *entities.ShopOwner
}

func validate(
	email valueobject.Email,
	address []AddressArg,
	firstName string,
	lastName string,
) error {
	if email == "" {
		return ErrEmailRequired
	}
	if firstName == "" {
		return ErrNameRequired
	}
	if lastName == "" {
		return ErrNameRequired
	}
	if len(address) == 0 {
		return ErrMinAddress
	}
	return nil
}

func NewUser(
	email valueobject.Email,
	password valueobject.Password,
	phoneNumber valueobject.PhoneNumber,
	addressArgs []AddressArg,
	firstName string,
	lastName string,
) (*User, error) {
	err := validate(email, addressArgs, firstName, lastName)
	if err != nil {
		return nil, err
	}
	id, err := valueobject.NewUserId(uuid.New())
	if err != nil {
		return nil, err
	}
	addresses := make([]*entities.Address, len(addressArgs))
	for idx, address := range addressArgs {
		addresses[idx] = entities.NewAddress(idx, address.Street, address.Town, address.City, address.Province)
	}
	// create customer info
	customer, err := entities.NewCustomer()
	if err != nil {
		return nil, err
	}
	return &User{
		Id:          *id,
		Email:       email,
		Password:    password,
		PhoneNumber: phoneNumber,
		Address:     addresses,
		FirstName:   firstName,
		LastName:    lastName,
		Customer:    customer,
	}, nil
}

func (u *User) SetCurrentRefreshToken(token string) {
	u.CurrentRefreshToken = token
}

func (u *User) ValidateAuth(password string) error {
	if !u.Password.Compare(password) {
		return ErrBadCredencial
	}
	return nil
}
func (u *User) ActivateShopOwner(args ActivateShopOwnerArg) {
	u.IsShopOwnerActive = true
	u.ShopOwner = entities.NewShopOwner(args.BussinessLincese)
}
