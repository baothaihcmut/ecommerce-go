package models

import (
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/enums"
	"github.com/google/uuid"
)

type UserContext struct {
	Id   uuid.UUID
	Role enums.Role
}
