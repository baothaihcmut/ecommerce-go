package models

type VariationValue struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}
