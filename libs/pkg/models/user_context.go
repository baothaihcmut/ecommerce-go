package models

import (
	"github.com/google/uuid"
)

type UserContext struct {
	Id   uuid.UUID
	Role Role
}
