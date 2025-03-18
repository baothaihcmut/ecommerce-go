package services

import (
	"context"
	"time"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/external"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


type AccessTokenClaims struct{
	jwt.RegisteredClaims
	UserId uuid.UUID `json:"user_id"`
	IsShopOwnerActive bool `json:"is_shop_owner_active"`
}

type RefreshTokenClaims struct{
	jwt.RegisteredClaims
	UserId uuid.UUID `json:"user_id"`
}


type JWTService struct {
	jwtConfig *config.JwtConfig
}



// GenerateAccessToken implements external.JwtService.
func (j *JWTService) GenerateAccessToken(ctx context.Context,payload external.AccessTokenPayload) (string, error) {
	claims := &AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.jwtConfig.AccessToken.Age) * time.Hour)),
		},
		UserId: payload.UserId,
		IsShopOwnerActive: payload.IsShopOwnerActive,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.jwtConfig.AccessToken.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateRefreshToken implements external.JwtService.
func (j *JWTService) GenerateRefreshToken(ctx context.Context,payload external.RefreshTokenPayload) (string, error) {
	claims := &RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.jwtConfig.RefreshToken.Age) * time.Hour)),
		},
		UserId: payload.UserId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.jwtConfig.RefreshToken.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func NewJWTService(jwtConfig *config.JwtConfig) external.JwtService {
	return &JWTService{
		jwtConfig: jwtConfig,
	}
}
