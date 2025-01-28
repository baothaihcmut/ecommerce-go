package services

import (
	"context"

	productitems "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/exceptions"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductItemCommandService struct {
	productItemRepo repositories.ProductItemCommandRepository
	productRepo     repositories.ProductCommandRepository
	mongoClient     *mongo.Client
}

func (p *ProductItemCommandService) CreateProductItem(ctx context.Context, command *commands.CreateProductItemCommand) (*results.CreateProductItemResult, error) {
	//check if product exist
	product, err := p.productRepo.FindById(ctx, command.ProductId)
	if err != nil {
		return nil, err
	}
	//check if variation not belong to product
	for _, variationValue := range command.VariationValues {
		if !product.CheckVariationBelongToProduct(variationValue.VariationId) {
			return nil, exceptions.ErrVariationOfItemNotBelongToProduct
		}
	}
	id := primitive.NewObjectID().Hex()
	productItemId := valueobjects.NewProductItemId(id)

	productItem, err := productitems.NewProductItem(
		productItemId,
		command.Sku,
		command.Price,
		command.Quantity,
		command.ProductId,
		command.VariationValues,
	)
	if err != nil {
		return nil, err
	}
	//save to db
	session, err := p.mongoClient.StartSession()
	if err != nil {
		return nil, err
	}
	session.StartTransaction()
	defer func() {
		if err != nil {
			session.AbortTransaction(ctx)
		}
	}()
	err = p.productItemRepo.Save(ctx, productItem, session)
	if err != nil {
		return nil, err
	}
	session.CommitTransaction(ctx)
	return &results.CreateProductItemResult{
		ProductItem: productItem,
	}, nil

}
func NewProductItemCommandService(
	productRepo repositories.ProductCommandRepository,
	productItemRepo repositories.ProductItemCommandRepository,
	mongoClient *mongo.Client,
) handlers.ProductItemCommandHandler {
	return &ProductItemCommandService{
		productRepo:     productRepo,
		productItemRepo: productItemRepo,
		mongoClient:     mongoClient,
	}
}
