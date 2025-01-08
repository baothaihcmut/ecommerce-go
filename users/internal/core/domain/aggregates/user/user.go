package user

import (
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
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

type User struct {
	Id                  valueobject.UserId
	Email               valueobject.Email
	Password            valueobject.Password
	PhoneNumber         valueobject.PhoneNumber
	Address             []valueobject.Address
	Role                enums.Role
	FirstName           string
	LastName            string
	CurrentRefreshToken *valueobject.Token
	Customer            *entities.Customer
	ShopOwner           *entities.ShopOwner
}

func validate(
	email valueobject.Email,
	address []valueobject.Address,
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

func NewCustomer(
	email valueobject.Email,
	password valueobject.Password,
	phoneNumber valueobject.PhoneNumber,
	address []valueobject.Address,
	firstName string,
	lastName string,
) (*User, error) {

	err := validate(email, address, firstName, lastName)
	if err != nil {
		return nil, err
	}
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
		Id:                  *id,
		Email:               email,
		Password:            password,
		CurrentRefreshToken: nil,
		PhoneNumber:         phoneNumber,
		Address:             address,
		Role:                enums.CUSTOMER,
		FirstName:           firstName,
		LastName:            lastName,
		Customer:            customer,
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
	if err != nil {
		return nil, err
	}
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
