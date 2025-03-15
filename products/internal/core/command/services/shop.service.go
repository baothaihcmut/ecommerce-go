package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/models"
)

type ShopService interface {
	FindShopById(ctx context.Context, id string) (*models.Shop, error)
}
