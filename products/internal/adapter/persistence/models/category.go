package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id               primitive.ObjectID `bson:"_id"`
	Name             string             `bson:"name"`
	ParentCategoryId []string           `bson:"parent_category_ids"`
}
