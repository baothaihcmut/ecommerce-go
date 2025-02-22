package external

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/models"
	"github.com/google/uuid"
)

type GenerateAccessTokenArg struct {
	UserId            uuid.UUID
	IsShopOwnerActive bool
}
type GenerateRefreshTokenArg struct {
	UserId uuid.UUID
}
type JwtService interface {
	DecodeAccessToken(context.Context, string) (models.AccessTokenSub, error)
	DecodeRefreshToken(context.Context, string) (models.RefreshTokenSub, error)
	GenerateAccessToken(context.Context, GenerateAccessTokenArg) (string, error)
	GenerateRefreshToken(context.Context, GenerateRefreshTokenArg) (string, error)
}
