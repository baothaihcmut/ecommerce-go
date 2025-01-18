package command

import (
	"context"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrParentCategoryNotExist = errors.New("parent category not exist")
)

type CategoryCommandService struct {
	categoryRepo repositories.CategoryCommandRepository
	mongoClient  *mongo.Client
}

func NewCategoryCommandService(repo repositories.CategoryCommandRepository) handlers.CategoryCommandHandler {
	return &CategoryCommandService{}
}

// CreateCategory implements handlers.CategoryCommandHandler.
func (c *CategoryCommandService) CreateCategory(ctx context.Context, command *commands.CreateCategoryCommand) (*results.CreateCategoryResult, error) {
	//if category is sub category check parent category exist
	parentCategoryIds := make([]valueobjects.CategoryId, len(command.ParentCategoryIds))
	for idx, val := range command.ParentCategoryIds {
		parentCategoryId := valueobjects.NewCategoryId(val)
		parentCategory, err := c.categoryRepo.FindCategoryById(ctx, parentCategoryId)
		if err != nil {
			return nil, err
		}
		if parentCategory == nil {
			return nil, ErrParentCategoryNotExist
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
	session, err := c.mongoClient.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer func(cause error) {
		if cause != nil {
			session.AbortTransaction(ctx)
		}
	}(err)
	c.categoryRepo.Save(ctx, category, session)
	res := &results.CreateCategoryResult{
		Id:               category.Id,
		Name:             category.Name,
		ParentCategoryId: parentCategoryIds,
	}
	return res, nil
}
