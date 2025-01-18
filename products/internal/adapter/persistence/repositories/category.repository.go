package repositories

import (
	"context"
	"fmt"
	"sync"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/filter"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/sort"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/persistence/models"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories/value_objects"
	commandRepository "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	queryRepository "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/outbound/repositories"
	categoryProjections "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/projections/categories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

type MongoCategoryRepository struct {
	collection *mongo.Collection
}

func decodeAnyValue[T protoreflect.ProtoMessage](value *anypb.Any) (protoreflect.Value, error) {
	var res T
	value.UnmarshalTo(res)
	return res.ProtoReflect().Get(), nil
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
func NewMongoCategoryCommandRepository(collection *mongo.Collection) commandRepository.CategoryCommandRepository {
	return &MongoCategoryRepository{
		collection: collection,
	}
}

func NewMongoCategoryQueryRepository(collection *mongo.Collection) queryRepository.CategoryQueryRepository {
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
	sessionCtx := mongo.NewSessionContext(ctx, session)
	opts := options.Update().SetUpsert(true)
	_, err = m.collection.UpdateOne(sessionCtx, bson.M{"_id": id}, bson.M{"$set": categoryModel}, opts)
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
) (*pagination.PaginationResult[*categoryProjections.CategoryProjection], error) {
	filterMongo := bson.M{}
	for _, filter := range filters {
		fmt.Println(filter.Field)
		fmt.Println(filter.Value)
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
