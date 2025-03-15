package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/persistence/models"
	categoryValueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories/value_objects"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products/entities"
	productValueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products/value_objects"
	valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products/value_objects"
	commonValueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/common/value_objects"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/outbound/repositories"
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
	images := make([]*entities.ProductImage, len(model.Images))
	for idx, image := range model.Images {
		urlSplit := strings.Split(image.Url, "/")
		images[idx] = entities.NewProductImage(
			valueobjects.NewProductImageId(productId, commonValueobjects.NewImageLink(urlSplit[0], urlSplit[1])),
			image.StorageProvider,
			image.Size,
			image.Type,
			image.Width,
			image.Height,
		)
	}
	return &products.Product{
		Id:          productId,
		Description: model.Description,
		Name:        model.Name,
		Unit:        model.Unit,
		CategoryIds: categoriesIds,
		Variations:  variations,
		Images:      images,
	}
}
func (m *MongoProductRepository) FindById(ctx context.Context, productId productValueobjects.ProductId) (*products.Product, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, m.tracer, "Product.FindById: database", nil)
	defer tracing.EndSpan(span, err, nil)
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
	var err error
	ctx, span := tracing.StartSpan(ctx, m.tracer, "Product.Save: database", nil)
	defer tracing.EndSpan(span, err, nil)
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
	productImages := make([]models.ProductImage, len(product.Images))
	for idx, val := range product.Images {
		productImages[idx] = models.ProductImage{
			Url:             fmt.Sprintf("%s/%s", val.Id.Url.Bucket, val.Id.Url.Key),
			Size:            val.Size,
			StorageProvider: val.StorageProvider,
			Type:            val.Type,
			Width:           val.Width,
			Height:          val.Height,
		}
	}
	productModel := &models.Product{
		Id:          id,
		Name:        product.Name,
		Description: product.Description,
		Unit:        product.Unit,
		ShopId:      string(product.ShopId),
		CategoryIds: categoryIds,
		Variations:  variations,
		Images:      productImages,
	}
	sessionCtx := mongo.NewSessionContext(ctx, session)
	opts := options.Update().SetUpsert(true)
	_, err = m.collection.UpdateOne(sessionCtx, bson.M{"_id": id}, bson.M{"$set": productModel}, opts)
	if err != nil {
		return err
	}
	return nil
}
