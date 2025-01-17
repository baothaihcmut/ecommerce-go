package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Unit        string             `bson:"unit"`
	ShopId      string             `bson:"shop_id"`
	CategoryIds []string           `bson:"category_ids"`
	Variations  []string           `bson:"variations"`
}
