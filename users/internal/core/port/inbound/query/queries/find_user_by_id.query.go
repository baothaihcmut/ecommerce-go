package queries

import "github.com/google/uuid"

type FindUserByIdQuery struct {
	Id uuid.UUID
}
