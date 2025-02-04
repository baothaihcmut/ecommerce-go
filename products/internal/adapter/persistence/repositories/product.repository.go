package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/persistence/models"
	categoryValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/entities"
	productValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/trace"
)

type MongoProductRepository struct {
	collection *mongo.Collection
	tracer     trace.Tracer
}

// FindById implements repositories.ProductCommandRepository.
func toProductDomain(model *models.Product) *products.Product {
	productId := productValueobjects.NewProductId(model.Id.Hex())
	categoriesIds := make([]categoryValueobjects.CategoryId, len(model.CategoryIds))
	for idx, categoryId := range model.CategoryIds {
		categoriesIds[idx] = categoryValueobjects.NewCategoryId(categoryId)
	}
	variations := make([]*entities.Variation, len(model.Variations))
	for idx, variation := range model.Variations {
		variations[idx] = entities.NewVariation(productValueobjects.NewVariationId(productId, variation))
	}
	return &products.Product{
		Id:          productId,
		Description: model.Description,
		Name:        model.Name,
		Unit:        model.Unit,
		CategoryIds: categoriesIds,
		Variations:  variations,
	}
}
func (m *MongoProductRepository) FindById(ctx context.Context, productId productValueobjects.ProductId) (*products.Product, error) {
	ctx, span := m.tracer.Start(ctx, "Product.FindById: database")
	defer span.End()
	id, err := primitive.ObjectIDFromHex(string(productId))
	if err != nil {
		return nil, err
	}
	productModel := models.Product{}
	err = m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&productModel)
	if err != nil {
		return nil, err
	}
	return toProductDomain(&productModel), nil

}

func NewMongoProductRepository(collection *mongo.Collection, tracer trace.Tracer) repositories.ProductCommandRepository {
	return &MongoProductRepository{
		collection: collection,
		tracer:     tracer,
	}
}

func (m *MongoProductRepository) Save(ctx context.Context, product *products.Product, session mongo.Session) error {
	ctx, span := m.tracer.Start(ctx, "Product.Save: database")
	defer span.End()
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
	_, err = m.collection.UpdateOne(sessionCtx, bson.M{"_id": id}, bson.M{"$set": productModel}, opts)
	if err != nil {
		return err
	}
	return nil
}
