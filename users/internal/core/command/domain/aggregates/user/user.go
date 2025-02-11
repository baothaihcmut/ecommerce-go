package user

import (
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
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

type User struct {
	Id                  valueobject.UserId
	Email               valueobject.Email
	Password            valueobject.Password
	PhoneNumber         valueobject.PhoneNumber
	Address             []*entities.Address
	Role                enums.Role
	FirstName           string
	LastName            string
	CurrentRefreshToken *valueobject.Token
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

func newUser(
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
	return &User{
		Id:                  *id,
		Email:               email,
		Password:            password,
		CurrentRefreshToken: nil,
		PhoneNumber:         phoneNumber,
		Address:             addresses,
		Role:                enums.CUSTOMER,
		FirstName:           firstName,
		LastName:            lastName,
	}, nil
}

func NewCustomer(
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
	user, err := newUser(email, password, phoneNumber, addressArgs, firstName, lastName)
	if err != nil {
		return nil, err
	}
	customer, err := entities.NewCustomer()
	user.Customer = customer
	// Create and return a new User instance
	return user, nil
}

func NewShopOwner(
	email valueobject.Email,
	passwod valueobject.Password,
	phoneNumber valueobject.PhoneNumber,
	addressArgs []AddressArg,
	firstName string,
	lastName string,
	bussinessLincese string,
) (*User, error) {
	err := validate(email, addressArgs, firstName, lastName)
	if err != nil {
		return nil, err
	}
	user, err := newUser(email, passwod, phoneNumber, addressArgs, firstName, lastName)
	if err != nil {
		return nil, err
	}

	shopOwner := entities.NewShopOwner(bussinessLincese)
	user.ShopOwner = shopOwner
	return user, nil
}

func (u *User) SetCurrentRefreshToken(token valueobject.Token) error {
	// check if refesh token is refresh token
	if token.TokenType != enums.REFRESH_TOKEN {
		return ErrInvalidTokenType
	}

	// if old refresh token != null compare
	if u.CurrentRefreshToken != nil && !u.CurrentRefreshToken.IsEqual(token) {
		return ErrMisMatchRefreshToken
	}

	u.CurrentRefreshToken = &token
	return nil
}

func (u *User) ValidateAuth(password string) error {
	if !u.Password.Compare(password) {
		return ErrBadCredencial
	}
	return nil
}
