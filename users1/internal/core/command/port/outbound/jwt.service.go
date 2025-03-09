package outbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/models"
	"github.com/google/uuid"
)

type GenerateTokenArg struct {
	UserId uuid.UUID
	Role   enums.Role
}

type JwtService interface {
	DecodeAccessToken(context.Context, string) (models.AccessTokenSub, error)
	DecodeRefreshToken(context.Context, string) (models.RefreshTokenSub, error)
	GenerateAccessToken(context.Context, GenerateTokenArg) (string, error)
	GenerateRefreshToken(context.Context, GenerateTokenArg) (string, error)
}
