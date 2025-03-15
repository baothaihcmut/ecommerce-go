package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/persistence/models"
	productitems "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/product_items"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/trace"
)

type MongoProductItemRepository struct {
	collection *mongo.Collection
	tracer     trace.Tracer
}

func NewMongoProductItemRepository(collection *mongo.Collection, tracer trace.Tracer) repositories.ProductItemCommandRepository {
	return &MongoProductItemRepository{
		collection: collection,
		tracer:     tracer,
	}
}

func (m *MongoProductItemRepository) Save(ctx context.Context, productItem *productitems.ProductItem, session mongo.Session) error {
	var err error
	ctx, span := tracing.StartSpan(ctx, m.tracer, "ProductItem.Save: database", nil)
	defer tracing.EndSpan(span, err, nil)
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
