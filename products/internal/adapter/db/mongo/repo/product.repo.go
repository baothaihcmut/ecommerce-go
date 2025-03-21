package repo

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/domain/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProductRepo struct {
	collection *mongo.Collection
}

func (m *MongoProductRepo) CreateProduct(ctx context.Context, p *entities.Product) error {
	_, err := m.collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}
func NewMongoProductRepo(db *mongo.Database) *MongoProductRepo {
	return &MongoProductRepo{
		collection: db.Collection("products"),
	}
}
