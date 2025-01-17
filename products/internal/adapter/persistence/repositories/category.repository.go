package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/persistence/models"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/categories"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCategoryRepository struct {
	collection *mongo.Collection
}

func NewMongoCategoryRepository(collection *mongo.Collection) repositories.CategoryRepository {
	return &MongoCategoryRepository{
		collection: collection,
	}
}

func (m *MongoCategoryRepository) Save(ctx context.Context, category *categories.Category, session mongo.Session) error {
	id, err := primitive.ObjectIDFromHex(string(category.Id))
	if err != nil {
		return err
	}
	//map category to bson
	parentCategoryIds := make([]string, len(category.ParentCategoryId))
	for idx, cate := range category.ParentCategoryId {
		parentCategoryIds[idx] = string(cate)
	}
	categoryModel := &models.Category{
		Id:               id,
		Name:             category.Name,
		ParentCategoryId: parentCategoryIds,
	}
	//
	sessionCtx := mongo.NewSessionContext(ctx, session)
	opts := options.Update().SetUpsert(true)
	_, err = m.collection.UpdateOne(sessionCtx, bson.M{"_id": id}, categoryModel, opts)
	if err != nil {
		return err
	}
	return nil
}
