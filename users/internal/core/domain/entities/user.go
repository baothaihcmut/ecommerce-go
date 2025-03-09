package entities

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
)

type CreateAddessArg struct {
	Priority int
	Street   string
	Town     string
	City     string
	Province string
}

type User struct {
	Id                  uuid.UUID
	Email               string
	Password            string
	PhoneNumber         string
	Addresses           []*Address
	IsShopOwnerActive   bool
	FirstName           string
	LastName            string
	CurrentRefreshToken *string
	Customer            *Customer
	ShopOwner           *ShopOwner
}

func hashPassword(passwd string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hashed)
}

func NewUser(
	email string,
	password string,
	phoneNumber string,
	addressArgs []CreateAddessArg,
	firstName string,
	lastName string,
) *User {
	userId := uuid.New()
	user := &User{
		Id:          userId,
		Email:       email,
		Password:    hashPassword(password),
		PhoneNumber: phoneNumber,
		Addresses: lo.Map(addressArgs, func(item CreateAddessArg, _ int) *Address {
			return NewAddress(
				userId,
				item.Priority,
				item.Street,
				item.Town,
				item.City,
				item.Province,
			)
		}),
		FirstName:           firstName,
		LastName:            lastName,
		CurrentRefreshToken: nil,
		ShopOwner:           nil,
	}
	user.Customer = NewCustomer(user)
	return user
}
