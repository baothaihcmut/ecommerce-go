package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shop struct {
	ID          primitive.ObjectID
	ShopOwnerId uuid.UUID
}
