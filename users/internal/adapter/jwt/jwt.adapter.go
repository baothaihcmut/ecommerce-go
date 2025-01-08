package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/models"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/outbound"
	services "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/services/command"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtAdapter struct {
	jwtConfig *config.JwtConfig
}

// DecodeRefreshToken implements outbound.JwtPort.
func (j *JwtAdapter) DecodeRefreshToken(_ context.Context, token valueobject.Token) (models.RefreshTokenSub, error) {
	tokenDecode, err := jwt.ParseWithClaims(token.Value, &JwtAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.jwtConfig.RefreshTokenSecret), nil
	})
	if err != nil {
		return models.RefreshTokenSub{}, err
	}
	if claims, ok := tokenDecode.Claims.(*JwtAccessClaims); ok && tokenDecode.Valid {
		if claims.ExpiresAt != nil {
			expireTime := claims.ExpiresAt.Time
			if time.Now().After(expireTime) {
				return models.RefreshTokenSub{}, services.ErrTokenExpire
			}
			return models.RefreshTokenSub{
				Id: valueobject.UserId(claims.UserId),
			}, nil
		}
	}
	return models.RefreshTokenSub{}, services.ErrInvalidToken
}

func (j *JwtAdapter) GenerateAccessToken(_ context.Context, u *user.User) (valueobject.Token, error) {
	claims := &JwtAccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.jwtConfig.AccessTokenAge) * time.Hour)),
		},
		UserId: uuid.UUID(u.Id),
		Role:   string(u.Role),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(j.jwtConfig.AccessTokenSecret)
	if err != nil {
		return valueobject.Token{}, err
	}
	return valueobject.Token{
		Value:     tokenString,
		TokenType: enums.ACCESS_TOKEN,
	}, nil
}

// GenerateRefreshToken implements outbound.JwtPort.
func (j *JwtAdapter) GenerateRefreshToken(_ context.Context, u *user.User) (valueobject.Token, error) {
	claims := &JwtRefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.jwtConfig.RefreshTokenAge) * time.Hour)),
		},
		UserId: uuid.UUID(u.Id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(j.jwtConfig.RefreshTokenSecret)
	if err != nil {
		return valueobject.Token{}, err
	}
	return valueobject.Token{
		Value:     tokenString,
		TokenType: enums.REFRESH_TOKEN,
	}, nil
}

func (j *JwtAdapter) DecodeAccessToken(_ context.Context, token valueobject.Token) (models.AccessTokenSub, error) {
	tokenDecode, err := jwt.ParseWithClaims(token.Value, &JwtAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.jwtConfig.AccessTokenSecret), nil
	})
	if err != nil {
		return models.AccessTokenSub{}, err
	}
	if claims, ok := tokenDecode.Claims.(*JwtAccessClaims); ok && tokenDecode.Valid {
		if claims.ExpiresAt != nil {
			expireTime := claims.ExpiresAt.Time
			if time.Now().After(expireTime) {
				return models.AccessTokenSub{}, services.ErrTokenExpire
			}
			return models.AccessTokenSub{
				Id:   valueobject.UserId(claims.UserId),
				Role: enums.Role(claims.Role),
			}, nil
		}
	}
	return models.AccessTokenSub{}, services.ErrInvalidToken
}

func NewJwtAdapter(jwtCfg *config.JwtConfig) outbound.JwtPort {
	return &JwtAdapter{jwtConfig: jwtCfg}
}
