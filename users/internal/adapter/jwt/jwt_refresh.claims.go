package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtRefreshClaims struct {
	jwt.RegisteredClaims
	UserId uuid.UUID
}
