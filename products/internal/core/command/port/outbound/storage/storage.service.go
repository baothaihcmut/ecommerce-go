package storage

import (
	"context"

	valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/common/value_objects"
)

type GetPresignUrlMethod string

const (
	GET GetPresignUrlMethod = "get"
	PUT GetPresignUrlMethod = "put"
)

type GetPresignUrlArgs struct {
	Link   valueobjects.ImageLink
	Method GetPresignUrlMethod
}

type StorageService interface {
	GetPresignUrl(context.Context, GetPresignUrlArgs) (string, error)
}
