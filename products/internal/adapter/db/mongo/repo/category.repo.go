package repo

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCategoryRepo struct {
	collection *mongo.Collection
}

func (m *MongoCategoryRepo) FindCategoryById(ctx context.Context, id primitive.ObjectID) (*entities.Category, error) {
	var res entities.Category
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}
func NewMongoCategoryRepo(db *mongo.Database) *MongoCategoryRepo {
	return &MongoCategoryRepo{
		collection: db.Collection("categories"),
	}
}
