package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/persistence/models"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(collection *mongo.Collection) repositories.ProductCommandRepository {
	return &MongoProductRepository{
		collection: collection,
	}
}

func (m *MongoProductRepository) Save(ctx context.Context, product *products.Product, session mongo.Session) error {
	id, err := primitive.ObjectIDFromHex(string(product.Id))
	if err != nil {
		return err
	}
	//map to model
	categoryIds := make([]string, len(product.CategoryIds))
	for idx, val := range product.CategoryIds {
		categoryIds[idx] = string(val)
	}
	variations := make([]string, len(product.Variations))
	for idx, val := range product.Variations {
		variations[idx] = val.Id.Name
	}
	productModel := &models.Product{
		Id:          id,
		Name:        product.Name,
		Description: product.Description,
		Unit:        product.Unit,
		ShopId:      string(product.ShopId),
		CategoryIds: categoryIds,
		Variations:  variations,
	}
	sessionCtx := mongo.NewSessionContext(ctx, session)
	opts := options.Update().SetUpsert(true)
	_, err = m.collection.UpdateOne(sessionCtx, bson.M{"_id": id}, productModel, opts)
	if err != nil {
		return err
	}
	return nil
}
