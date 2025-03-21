package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductItem struct {
	ID              primitive.ObjectID `bson:"_id"`
	ProductId       primitive.ObjectID `bson:"product_id"`
	Quantity        int                `bson:"quantity"`
	Price           int                `bson:"price"`
	Images          []string           `bson:"images"`
	VariationValues []VariationValue   `bson:"variation_values"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}

type VariationValue struct {
	Variation string `bson:"variation"`
	Value     string `bson:"value"`
}
