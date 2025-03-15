package entities

import (
	valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products/value_objects"
)

type ProductImage struct {
	Id              valueobjects.ProductImageId
	StorageProvider string
	Size            int
	Type            string
	Width           int
	Height          int
}

func NewProductImage(
	Id valueobjects.ProductImageId,
	StorageProvider string,
	Size int,
	Type string,
	Width int,
	Height int,
) *ProductImage {
	return &ProductImage{
		Id:              Id,
		Size:            Size,
		Width:           Width,
		Height:          Height,
		StorageProvider: StorageProvider,
		Type:            Type,
	}
}
