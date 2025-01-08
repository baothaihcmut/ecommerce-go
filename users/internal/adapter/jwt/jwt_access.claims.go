package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtAccessClaims struct {
	jwt.RegisteredClaims
	UserId uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
}
