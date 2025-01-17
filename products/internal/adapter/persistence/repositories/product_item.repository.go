package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/persistence/models"
	productitems "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/product_items"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProductItemRepository struct {
	collection *mongo.Collection
}

func NewMongoProductItemRepository(collection *mongo.Collection) repositories.ProductItemRepository {
	return &MongoProductItemRepository{
		collection: collection,
	}
}

func (m *MongoProductItemRepository) Save(ctx context.Context, productItem *productitems.ProductItem, session mongo.Session) error {
	id, err := primitive.ObjectIDFromHex(string(productItem.Id))
	if err != nil {
		return err
	}

	// Map variation values to a suitable structure
	variationValues := make([]models.VariationValue, len(productItem.VariationValues))
	for idx, variation := range productItem.VariationValues {
		variationValues[idx] = models.VariationValue{
			Name:  variation.VariationId.Name,
			Value: variation.Value,
		}
	}

	// Map productItem to bson
	productItemModel := &models.ProductItem{
		Id:              id,
		Sku:             string(productItem.Sku),
		Price:           float64(productItem.Price),
		Quantity:        int(productItem.Quantity),
		ProductId:       string(productItem.ProductId),
		VariationValues: variationValues,
	}

	sessionCtx := mongo.NewSessionContext(ctx, session)
	opts := options.Update().SetUpsert(true)

	// Use UpdateOne with upsert option
	_, err = m.collection.UpdateOne(
		sessionCtx,
		bson.M{"_id": id},
		bson.M{"$set": productItemModel},
		opts,
	)
	if err != nil {
		return err
	}
	return nil
}
