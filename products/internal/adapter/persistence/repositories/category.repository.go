package repositories

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/filter"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/sort"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/persistence/models"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories"
	valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories/value_objects"
	commandRepository "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/outbound/repositories"
	queryRepository "github.com/baothaihcmut/Ecommerce-go/products/internal/core/query/port/outbound/repositories"
	categoryProjections "github.com/baothaihcmut/Ecommerce-go/products/internal/core/query/projections/categories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/trace"
)

type MongoCategoryRepository struct {
	collection *mongo.Collection
	tracer     trace.Tracer
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
func toCategoryProjection(model *models.Category) *categoryProjections.CategoryProjection {
	return &categoryProjections.CategoryProjection{
		Id:                model.Id.Hex(),
		Name:              model.Name,
		ParentCategoryIds: model.ParentCategoryId,
	}
}

func handleRoutineError(ctx context.Context, cancel context.CancelFunc, errCh chan error, err error) {
	select {
	case <-ctx.Done():
		break
	default:
		cancel()
		errCh <- err
	}
}
func NewMongoCategoryCommandRepository(collection *mongo.Collection, tracer trace.Tracer) commandRepository.CategoryCommandRepository {
	return &MongoCategoryRepository{
		collection: collection,
		tracer:     tracer,
	}
}

func NewMongoCategoryQueryRepository(collection *mongo.Collection, tracer trace.Tracer) queryRepository.CategoryQueryRepository {
	return &MongoCategoryRepository{
		collection: collection,
		tracer:     tracer,
	}
}

func (m *MongoCategoryRepository) Save(ctx context.Context, category *categories.Category, session mongo.Session) error {
	var err error
	ctx, span := tracing.StartSpan(ctx, m.tracer, "Category.Save: database", map[string]interface{}{
		"category_id": string(category.Id),
	})
	defer tracing.EndSpan(span, err, nil)
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
	sessionCtx := mongo.NewSessionContext(ctx, session)
	opts := options.Update().SetUpsert(true)
	_, err = m.collection.UpdateOne(sessionCtx, bson.M{"_id": id}, bson.M{"$set": categoryModel}, opts)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoCategoryRepository) BulkSave(ctx context.Context, categories []*categories.Category, session mongo.Session) error {
	var err error
	ctx, span := tracing.StartSpan(ctx, m.tracer, "Category.BulkSave: database", nil)
	defer tracing.EndSpan(span, err, nil)
	bulkOptions := make([]mongo.WriteModel, len(categories))
	for idx, val := range categories {
		id, err := primitive.ObjectIDFromHex(string(val.Id))
		if err != nil {
			return err
		}
		parentCategoryIds := make([]string, len(val.ParentCategoryId))
		for idx, cate := range val.ParentCategoryId {
			parentCategoryIds[idx] = string(cate)
		}
		categoryModel := &models.Category{
			Id:               id,
			Name:             val.Name,
			ParentCategoryId: parentCategoryIds,
		}
		upsertOpt := mongo.NewUpdateManyModel().
			SetFilter(bson.M{
				"_id": id,
			}).
			SetUpdate(bson.M{
				"$set": categoryModel,
			}).
			SetUpsert(true)
		bulkOptions[idx] = upsertOpt
	}
	sessionCtx := mongo.NewSessionContext(ctx, session)
	_, err = m.collection.BulkWrite(sessionCtx, bulkOptions)
	if err != nil {
		return err
	}
	return nil

}

func (m *MongoCategoryRepository) FindCategoryById(ctx context.Context, categoryId valueobjects.CategoryId) (*categories.Category, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, m.tracer, "Category.FindById: database", nil)
	defer tracing.EndSpan(span, err, nil)
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
) (*pagination.PaginationResult[*categoryProjections.CategoryProjection], error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, m.tracer, "Category.FindAllCategory: database", nil)
	defer tracing.EndSpan(span, err, nil)
	filterMongo := bson.M{}
	for _, filter := range filters {
		filterMongo[filter.Field] = filter.Value
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
	//channel for data
	dataCh := make(chan []*categoryProjections.CategoryProjection, 1)
	//chanel for count
	countCh := make(chan int, 1)
	//chanel for error
	errCh := make(chan error, 1)
	//context for cancel when have error
	ctx, cancel := context.WithCancel(ctx)
	//wait group
	wg := &sync.WaitGroup{}

	//routine for query data
	wg.Add(1)
	go func() {
		defer wg.Done()
		//for filter
		cursor, err := m.collection.Find(ctx, filterMongo, findOptions)
		if err != nil {
			handleRoutineError(ctx, cancel, errCh, err)
			return
		}
		defer cursor.Close(ctx)
		var categorieMongo []models.Category
		if err = cursor.All(ctx, &categorieMongo); err != nil {
			handleRoutineError(ctx, cancel, errCh, err)
			return
		}
		res := make([]*categoryProjections.CategoryProjection, len(categorieMongo))
		for idx, model := range categorieMongo {
			res[idx] = toCategoryProjection(&model)
		}
		dataCh <- res
	}()
	//routine for count
	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := m.collection.CountDocuments(ctx, filterMongo)
		if err != nil {
			handleRoutineError(ctx, cancel, errCh, err)
			return
		}
		countCh <- int(count)
	}()
	//wait all routine done
	wg.Wait()
	//if have error return err
	select {
	case err := <-errCh:
		return nil, err
	default:
	}
	data := <-dataCh
	count := <-countCh
	return &pagination.PaginationResult[*categoryProjections.CategoryProjection]{
		Data: data,
		Pagination: pagination.Pagination{
			CurrentPage: paginate.Page,
			PageSize:    paginate.Size,
			TotalPage:   count / paginate.Size,
			TotalItem:   count,
		},
	}, nil
}

func (m *MongoCategoryRepository) FindAllSubCategory(ctx context.Context, categoryId string) ([]*categoryProjections.CategoryProjection, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, m.tracer, "Category.FindAllCategory: database", nil)
	defer tracing.EndSpan(span, err, nil)
	filter := bson.M{
		"parent_category_ids": bson.M{
			"$in": bson.A{categoryId},
		},
	}
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var categoryModels []*models.Category
	if err := cursor.All(ctx, &categoryModels); err != nil {
		return nil, err
	}
	mapResWg := &sync.WaitGroup{}
	res := make([]*categoryProjections.CategoryProjection, len(categoryModels))
	for idx, category := range categoryModels {
		mapResWg.Add(1)
		go func() {
			defer mapResWg.Done()
			res[idx] = toCategoryProjection(category)
		}()
	}
	mapResWg.Wait()
	return res, nil
}
