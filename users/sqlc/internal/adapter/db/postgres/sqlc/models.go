// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Address struct {
	Priority int32
	Street   string
	Town     string
	City     string
	Province string
	UserID   pgtype.UUID
}

type Admin struct {
	ID                  pgtype.UUID
	Email               string
	Password            string
	PhoneNumber         string
	FirstName           string
	LastName            string
	CurrentRefreshToken pgtype.Text
}

type Customer struct {
	UserID     pgtype.UUID
	LoyalPoint pgtype.Int4
}

type ShopOwner struct {
	UserID           pgtype.UUID
	BussinessLicense pgtype.Text
}

type User struct {
	ID                  pgtype.UUID
	Email               string
	Password            string
	IsShopOwnerActive   bool
	PhoneNumber         string
	FirstName           string
	LastName            string
	CurrentRefreshToken pgtype.Text
}
