package domain

import (
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain/enums"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain/value_object"
	"github.com/google/uuid"
)

var (
	ErrMisMatchRefreshToken = errors.New("Mismatch refrresh token")
	InValidTokenType        = errors.New("Invalid token type")
)

type User struct {
	Id                  valueobject.UserId
	Email               valueobject.Email
	Role                enums.Role
	CurrentRefreshToken *valueobject.Token
}

func NewUser(email valueobject.Email, role enums.Role, currentRefreshToken *valueobject.Token) (*User, error) {
	if currentRefreshToken.TokenType != enums.REFRESH_TOKEN {
		return nil, InValidTokenType
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &User{
		Id:                  valueobject.UserId(id),
		Email:               email,
		CurrentRefreshToken: currentRefreshToken,
		Role:                role,
	}, nil
}
func (u *User) CompareRefreshToken(token valueobject.Token) error {
	return u.CurrentRefreshToken.IsEqual(token)

}
