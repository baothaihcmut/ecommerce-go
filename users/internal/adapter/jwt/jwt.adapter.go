package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/errors"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/models"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/outbound"
	services "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/services/command"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtAdapter struct {
	accessTokenSecret  string
	accessTokenAge     int
	refreshTokenSecret string
	refreshTokenAge    int
}

// DecodeRefreshToken implements outbound.JwtPort.
func (j *JwtAdapter) DecodeRefreshToken(_ context.Context, token string) (models.RefreshTokenSub, error) {
	tokenDecode, err := jwt.ParseWithClaims(token, &JwtAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.refreshTokenSecret), nil
	})
	if err != nil {
		return models.RefreshTokenSub{}, errors.NewError(err, errors.CaptureStackTrace())
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

func (j *JwtAdapter) GenerateAccessToken(_ context.Context, args outbound.GenerateTokenArg) (string, error) {
	claims := &JwtAccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.accessTokenAge) * time.Hour)),
		},
		UserId: uuid.UUID(args.UserId),
		Role:   string(args.Role),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.accessTokenSecret))
	if err != nil {
		return "", errors.NewError(err, errors.CaptureStackTrace())
	}
	return tokenString, nil
}

// GenerateRefreshToken implements outbound.JwtPort.
func (j *JwtAdapter) GenerateRefreshToken(_ context.Context, args outbound.GenerateTokenArg) (string, error) {
	claims := &JwtRefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.refreshTokenAge) * time.Hour)),
		},
		UserId: uuid.UUID(args.UserId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.refreshTokenSecret))
	if err != nil {
		return "", errors.NewError(err, errors.CaptureStackTrace())
	}
	return tokenString, nil
}

func (j *JwtAdapter) DecodeAccessToken(_ context.Context, token string) (models.AccessTokenSub, error) {

	tokenDecode, err := jwt.ParseWithClaims(token, &JwtAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.accessTokenSecret), nil
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

func NewJwtService(jwtCfg *config.JwtConfig) outbound.JwtService {
	return &JwtAdapter{
		accessTokenSecret:  jwtCfg.AccessTokenSecret,
		accessTokenAge:     jwtCfg.AccessTokenAge,
		refreshTokenSecret: jwtCfg.RefreshTokenSecret,
		refreshTokenAge:    jwtCfg.RefreshTokenAge,
	}
}

func NewAdminJwtService(adminCfg *config.AdminConfig) outbound.JwtService {
	return &JwtAdapter{
		accessTokenSecret:  adminCfg.AccessTokenSecret,
		accessTokenAge:     adminCfg.AccessTokenAge,
		refreshTokenSecret: adminCfg.RefreshTokenSecret,
		refreshTokenAge:    adminCfg.RefreshTokenAge,
	}
}
