package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtAccessClaims struct {
	jwt.RegisteredClaims
	UserId            uuid.UUID `json:"userId"`
	IsShopOwnerActive bool      `json:"isShopOwnerActive"`
}
