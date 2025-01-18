package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/filter"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/sort"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/persistence/models"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCategoryRepository struct {
	collection *mongo.Collection
}

func toCategoryDomain(model *models.Category) *categories.Category {
	parentCategoryIds := make([]valueobjects.CategoryId, len(model.ParentCategoryId))
	for idx, val := range model.ParentCategoryId {
		parentCategoryIds[idx] = valueobjects.CategoryId(val)
	}
	return &categories.Category{
		Id:               valueobjects.NewCategoryId(model.Id.Hex()),
		Name:             model.Name,
		ParentCategoryId: parentCategoryIds,
	}
}
func NewMongoCategoryRepository(collection *mongo.Collection) repositories.CategoryCommandRepository {
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

func (m *MongoCategoryRepository) FindCategoryById(ctx context.Context, categoryId valueobjects.CategoryId) (*categories.Category, error) {
	id, err := primitive.ObjectIDFromHex(string(categoryId))
	if err != nil {
		return nil, err
	}
	categoryModel := models.Category{}

	err = m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&categoryModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	//map to domain
	return toCategoryDomain(&categoryModel), nil
}
func (m *MongoCategoryRepository) FindAllCategory(
	ctx context.Context,
	filters []filter.FilterParam,
	sorts []sort.SortParam,
	paginate pagination.PaginationParam,

) ([]*categories.Category, error) {
	//for filter
	filterArr := bson.D{}
	for idx, filter := range filters {
		filterArr[idx] = bson.E{Key: filter.Field, Value: filter.Value}
	}
	filterMongo := bson.M{
		"$and": filterArr,
	}
	//for sort
	sortMongo := bson.M{}
	for _, sortParam := range sorts {
		if sortParam.Direction == sort.ASC {
			sortMongo[sortParam.Field] = 1
		} else {
			sortMongo[sortParam.Field] = -1
		}
	}
	limit := paginate.Size
	skip := (paginate.Page - 1) * limit
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(sortMongo)

	cursor, err := m.collection.Find(ctx, filterMongo, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var categorieMongo []models.Category
	if err = cursor.All(ctx, &categorieMongo); err != nil {
		return nil, err
	}
	res := make([]*categories.Category, len(categorieMongo))
	for idx, model := range categorieMongo {
		res[idx] = toCategoryDomain(&model)
	}
	return res, nil

}
