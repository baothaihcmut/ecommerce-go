package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductItem struct {
	Id              primitive.ObjectID `bson:"_id"`
	Sku             string             `bson:"sku"`
	Price           float64            `bson:"price"`
	Quantity        int                `bson:"quantity"`
	ProductId       string             `bson:"product_id"`
	VariationValues []VariationValue   `bson:"variation_values"`
}
