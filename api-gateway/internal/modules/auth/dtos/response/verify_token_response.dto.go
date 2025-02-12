package response

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/models"
	"github.com/google/uuid"
)

type VerifyTokenResponse struct {
	Id   uuid.UUID
	Role models.Role
}
