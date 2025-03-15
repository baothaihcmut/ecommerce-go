package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/external"
)

type JWTService struct {
}

// GenerateAccessToken implements external.JwtService.
func (j *JWTService) GenerateAccessToken(context.Context, external.AccessTokenPayload) (string, error) {
	return "", nil
}

// GenerateRefreshToken implements external.JwtService.
func (j *JWTService) GenerateRefreshToken(context.Context, external.RefreshTokenPayload) (string, error) {
	return "", nil
}

func NewJWTService() external.JwtService {
	return &JWTService{}
}
