package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/mongo"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	productitems "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/exceptions"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type ProductItemCommandService struct {
	productItemRepo    repositories.ProductItemCommandRepository
	productRepo        repositories.ProductCommandRepository
	transactionService mongo.MongoTransactionService
	tracer             trace.Tracer
}

func (p *ProductItemCommandService) CreateProductItem(ctx context.Context, command *commands.CreateProductItemCommand) (*results.CreateProductItemResult, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, p.tracer, "ProductItem.Create: service", nil)
	defer tracing.EndSpan(span, err, nil)
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
	session, err := p.transactionService.BeginTransaction(ctx)
	defer func() {
		if err != nil {
			p.transactionService.RollbackTransaction(ctx, session)
		}
		p.transactionService.EndTransansaction(ctx, session)
	}()
	err = p.productItemRepo.Save(ctx, productItem, session)
	if err != nil {
		return nil, err
	}
	err = p.transactionService.CommitTransaction(ctx, session)
	if err != nil {
		return nil, err
	}
	return &results.CreateProductItemResult{
		ProductItem: productItem,
	}, nil

}
func NewProductItemCommandService(
	productRepo repositories.ProductCommandRepository,
	productItemRepo repositories.ProductItemCommandRepository,
	mongoClient mongo.MongoTransactionService,
	tracer trace.Tracer,
) handlers.ProductItemCommandHandler {
	return &ProductItemCommandService{
		productRepo:        productRepo,
		productItemRepo:    productItemRepo,
		transactionService: mongoClient,
		tracer:             tracer,
	}
}
