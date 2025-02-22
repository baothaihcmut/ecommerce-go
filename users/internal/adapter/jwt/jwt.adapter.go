package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/errors"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/models"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound/external"
	services "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/services"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type JwtAdapter struct {
	accessTokenSecret  string
	accessTokenAge     int
	refreshTokenSecret string
	refreshTokenAge    int
	tracer             trace.Tracer
}

// DecodeRefreshToken implements outbound.JwtPort.
func (j *JwtAdapter) DecodeRefreshToken(ctx context.Context, token string) (_ models.RefreshTokenSub, err error) {
	_, span := tracing.StartSpan(ctx, j.tracer, "Jwt.DecodeRefreshToken: service", nil)
	defer tracing.EndSpan(span, err, nil)
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

func (j *JwtAdapter) GenerateAccessToken(ctx context.Context, args external.GenerateAccessTokenArg) (_ string, err error) {
	ctx, span := tracing.StartSpan(ctx, j.tracer, "Jwt.GenerateAccessToken: service", nil)
	defer tracing.EndSpan(span, err, nil)
	claims := &JwtAccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.accessTokenAge) * time.Hour)),
		},
		UserId:            uuid.UUID(args.UserId),
		IsShopOwnerActive: args.IsShopOwnerActive,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.accessTokenSecret))
	if err != nil {
		return "", errors.NewError(err, errors.CaptureStackTrace())
	}
	return tokenString, nil
}

// GenerateRefreshToken implements outbound.JwtPort.
func (j *JwtAdapter) GenerateRefreshToken(ctx context.Context, args external.GenerateRefreshTokenArg) (_ string, err error) {
	ctx, span := tracing.StartSpan(ctx, j.tracer, "Jwt.GenerateRefreshToken: service", nil)
	defer tracing.EndSpan(span, err, nil)
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

func (j *JwtAdapter) DecodeAccessToken(ctx context.Context, token string) (_ models.AccessTokenSub, err error) {
	ctx, span := tracing.StartSpan(ctx, j.tracer, "Jwt.DecodeAcessToken: service", nil)
	defer tracing.EndSpan(span, err, nil)
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
				Id:                valueobject.UserId(claims.UserId),
				IsShopOwnerActive: claims.IsShopOwnerActive,
			}, nil
		}
	}
	return models.AccessTokenSub{}, services.ErrInvalidToken
}

func NewJwtService(jwtCfg *config.JwtConfig, tracer trace.Tracer) external.JwtService {
	return &JwtAdapter{
		accessTokenSecret:  jwtCfg.AccessTokenSecret,
		accessTokenAge:     jwtCfg.AccessTokenAge,
		refreshTokenSecret: jwtCfg.RefreshTokenSecret,
		refreshTokenAge:    jwtCfg.RefreshTokenAge,
		tracer:             tracer,
	}
}

func NewAdminJwtService(adminCfg *config.AdminConfig, tracer trace.Tracer) external.JwtService {
	return &JwtAdapter{
		accessTokenSecret:  adminCfg.AccessTokenSecret,
		accessTokenAge:     adminCfg.AccessTokenAge,
		refreshTokenSecret: adminCfg.RefreshTokenSecret,
		refreshTokenAge:    adminCfg.RefreshTokenAge,
		tracer:             tracer,
	}
}
