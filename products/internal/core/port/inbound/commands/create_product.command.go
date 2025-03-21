package commands

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateProductCommand struct {
	Name         string
	Description  string
	ShopId       primitive.ObjectID
	CategoryIds  []primitive.ObjectID
	NumOfImages  int
	HasThumbNail bool
	Variations   []string
}
