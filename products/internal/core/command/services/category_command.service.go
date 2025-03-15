package services

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/mongo"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories"
	valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories/value_objects"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/exceptions"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type CategoryCommandService struct {
	categoryRepo       repositories.CategoryCommandRepository
	transactionService mongo.MongoTransactionService
	tracer             trace.Tracer
}

func NewCategoryCommandService(repo repositories.CategoryCommandRepository, mongoClient mongo.MongoTransactionService, tracer trace.Tracer) handlers.CategoryCommandHandler {
	return &CategoryCommandService{
		categoryRepo:       repo,
		transactionService: mongoClient,
		tracer:             tracer,
	}
}

func (c *CategoryCommandService) checkCategoryExist(ctx context.Context, categoryId valueobjects.CategoryId) error {
	parentCategory, err := c.categoryRepo.FindCategoryById(ctx, categoryId)
	if err != nil {
		return err
	}
	if parentCategory == nil {
		return exceptions.ErrParentCategoryNotExist
	}
	return nil
}

func (c *CategoryCommandService) toCreateCategoryResult(category *categories.Category) *results.CreateCategoryResult {
	res := &results.CreateCategoryResult{
		Category: category,
	}
	return res
}

// CreateCategory implements handlers.CategoryCommandHandler.
func (c *CategoryCommandService) CreateCategory(ctx context.Context, command *commands.CreateCategoryCommand) (*results.CreateCategoryResult, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, c.tracer, "Category.Create: service", nil)
	defer tracing.EndSpan(span, err, nil)
	//if category is sub category check parent category exist
	parentCategoryIds := make([]valueobjects.CategoryId, len(command.ParentCategoryIds))
	for idx, val := range command.ParentCategoryIds {
		parentCategoryId := valueobjects.NewCategoryId(val)
		err = c.checkCategoryExist(ctx, parentCategoryId)
		if err != nil {
			return nil, err
		}
		parentCategoryIds[idx] = parentCategoryId
	}

	//create new category domain
	categoryId := valueobjects.NewCategoryId(primitive.NewObjectID().Hex())
	category := categories.NewCategory(
		categoryId,
		command.Name,
		parentCategoryIds,
	)
	//persist to db
	session, err := c.transactionService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = c.transactionService.RollbackTransaction(ctx, session)
		}
		c.transactionService.EndTransansaction(ctx, session)
	}()
	err = c.categoryRepo.Save(ctx, category, session)
	if err != nil {
		return nil, err
	}
	err = c.transactionService.CommitTransaction(ctx, session)
	if err != nil {
		return nil, err
	}
	return c.toCreateCategoryResult(category), nil
}

func (c *CategoryCommandService) BulkCreateCategories(ctx context.Context, command *commands.BulkCreateCategories) (*results.BulkCreateCategoriesResult, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, c.tracer, "Category.BulkCreate: service", nil)
	defer tracing.EndSpan(span, err, nil)
	//create set of all parent category need to check
	categorySet := make(map[string]struct{})
	//mutex lock for set
	categorySetLock := &sync.Mutex{}
	//waitgroup
	wg := &sync.WaitGroup{}
	categoryDomains := make([]*categories.Category, len(command.Categories))

	for idx, category := range command.Categories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parentCategoryIds := make([]valueobjects.CategoryId, len(category.ParentCategoryIds))
			for parenCategoryIdx, parentCategory := range category.ParentCategoryIds {
				wg.Add(1)
				go func() {
					defer wg.Done()
					parentCategoryIds[parenCategoryIdx] = valueobjects.NewCategoryId(parentCategory)
					categorySetLock.Lock()
					_, exist := categorySet[parentCategory]
					if !exist {
						categorySet[parentCategory] = struct{}{}
					}
					categorySetLock.Unlock()
				}()
			}
			//create new domain
			categoryDomains[idx] = categories.NewCategory(
				valueobjects.NewCategoryId(primitive.NewObjectID().Hex()),
				category.Name,
				parentCategoryIds,
			)
		}()
	}
	wg.Wait()
	//check category exist in set
	checkExistWg := &sync.WaitGroup{}
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for category := range categorySet {
		checkExistWg.Add(1)
		go func() {
			defer checkExistWg.Done()
			if err = c.checkCategoryExist(ctx, valueobjects.NewCategoryId(category)); err != nil {
				select {
				case <-ctx.Done():
				default:
					cancel()
					errCh <- err
				}
			}
		}()
	}
	checkExistWg.Wait()
	select {
	case err = <-errCh:
		return nil, err
	default:
	}
	//persist to db
	session, err := c.transactionService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = c.transactionService.RollbackTransaction(ctx, session)
		}
		c.transactionService.EndTransansaction(ctx, session)
	}()
	err = c.categoryRepo.BulkSave(ctx, categoryDomains, session)
	if err != nil {
		return nil, err
	}
	err = c.transactionService.CommitTransaction(ctx, session)
	if err != nil {
		return nil, err
	}
	//map domain to result
	mapResultWg := &sync.WaitGroup{}
	categoryResults := make([]*results.CreateCategoryResult, len(categoryDomains))
	for idx, category := range categoryDomains {
		mapResultWg.Add(1)
		go func() {
			defer mapResultWg.Done()
			categoryResults[idx] = c.toCreateCategoryResult(category)
		}()
	}
	mapResultWg.Wait()
	return &results.BulkCreateCategoriesResult{
		Categories: categoryResults,
	}, nil

}
