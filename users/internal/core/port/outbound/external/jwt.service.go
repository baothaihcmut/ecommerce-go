package external

import (
	"context"

	"github.com/google/uuid"
)

type AccessTokenPayload struct {
	UserId            uuid.UUID
	IsShopOwnerActive bool
}

type RefreshTokenPayload struct {
	UserId uuid.UUID
}

type JwtService interface {
	GenerateAccessToken(context.Context, AccessTokenPayload) (string, error)
	GenerateRefreshToken(context.Context, RefreshTokenPayload) (string, error)
}
