package models

import (
	"github.com/google/uuid"
)

type UserContext struct {
	UserId            uuid.UUID
	IsShopOwnerActive bool
}
