package valueobject

import "github.com/google/uuid"

type UserId uuid.UUID

func NewUserId(id uuid.UUID) (*UserId, error) {
	return (*UserId)(&id), nil
}
