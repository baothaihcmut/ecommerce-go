package services

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/mongo"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"
	commonValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/common/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/exceptions"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type ProductCommandService struct {
	productRepo        repositories.ProductCommandRepository
	categoryRepo       repositories.CategoryCommandRepository
	shopService        ShopService
	transactionService mongo.MongoTransactionService
	tracer             trace.Tracer
	s3Config           *config.S3Config
	storageService     storage.StorageService
}

func (p *ProductCommandService) checkContraints(ctx context.Context, product *commands.CreateProductCommand) error {
	//check categories and shop exist
	checkExistWg := &sync.WaitGroup{}
	errCh := make(chan error, 1)
	//context for cancel when 1 categories not exist
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, categoryId := range product.CategoryIds {
		checkExistWg.Add(1)
		go func() {
			defer checkExistWg.Done()
			categoryExist, err := p.categoryRepo.FindCategoryById(ctx, categoryId)
			if err != nil {
				if err == context.Canceled {
					return
				}
				cancel()
				errCh <- err
				return
			}
			if categoryExist == nil {
				errCh <- exceptions.ErrCategoryNotExist
			}
		}()
	}
	checkExistWg.Add(1)
	go func() {
		defer checkExistWg.Done()
		shop, err := p.shopService.FindShopById(ctx, string(product.ShopId))
		if err != nil {
			if err == context.Canceled {
				return
			}
			cancel()
			errCh <- err
			return
		}
		//if shop not exist
		if shop == nil {
			errCh <- exceptions.ErrShopNotExist
			return
		}
		//if shop not active
		if !shop.IsActive {
			errCh <- exceptions.ErrShopNotActive
			return
		}
	}()
	checkExistWg.Wait()
	//if have error return
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (p *ProductCommandService) CreateProduct(ctx context.Context, product *commands.CreateProductCommand) (*results.CreateProductResult, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, p.tracer, "Product.Create: endpoint", nil)
	defer tracing.EndSpan(span, err, nil)
	//check contraint
	err = p.checkContraints(ctx, product)
	if err != nil {
		return nil, err
	}
	id := primitive.NewObjectID().Hex()
	productId := valueobjects.NewProductId(id)
	//key for image

	imageArgs := make([]products.ProductImageArgs, len(product.Images))
	for idx, image := range product.Images {
		imageArgs[idx] = products.ProductImageArgs{
			StorageProvider: p.s3Config.StorageProvider,
			Size:            image.Size,
			Type:            image.Type,
			Width:           image.Width,
			Height:          image.Height,
			Url:             commonValueobjects.NewImageLink(p.s3Config.Bucket, id),
		}
	}

	productDomain, err := products.NewProduct(
		productId,
		product.Name,
		product.Description,
		product.Unit,
		product.ShopId,
		product.CategoryIds,
		product.Variations,
		imageArgs,
	)
	if err != nil {
		return nil, err
	}
	//save to db
	session, err := p.transactionService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			p.transactionService.RollbackTransaction(ctx, session)
		}
		p.transactionService.EndTransansaction(ctx, session)
	}()
	err = p.productRepo.Save(ctx, productDomain, session)
	if err != nil {
		return nil, err
	}
	err = p.transactionService.CommitTransaction(ctx, session)
	if err != nil {
		return nil, err
	}
	//get presign url for put
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for idx, image := range productDomain.Images {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := p.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArgs{
				Link:   image.Id.Url,
				Method: storage.PUT,
			})
			if err != nil {
				if err == context.Canceled {
					return
				}
				cancel()
				errCh <- err
				return
			}
			productDomain.Images[idx].Id.Url.Key = url
		}()
	}
	wg.Wait()
	select {
	case err = <-errCh:
		return nil, err
	default:
	}
	return &results.CreateProductResult{
		Product: productDomain,
	}, nil
}

func (p *ProductCommandService) UpdateProduct(ctx context.Context, command *commands.UpdateProductCommand) (*results.UpdateProductResult, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, p.tracer, "Product.Update: service", nil)
	defer tracing.EndSpan(span, err, nil)
	//find by id
	product, err := p.productRepo.FindById(ctx, command.Id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, exceptions.ErrProductNotExist
	}
	//update
	if command.Name != nil {
		product.Name = *command.Name
	}
	if command.Description != nil {
		product.Description = *command.Description
	}
	if command.Unit != nil {
		product.Unit = *command.Unit
	}
	//save to db
	session, err := p.transactionService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			p.transactionService.RollbackTransaction(ctx, session)
		}
		p.transactionService.EndTransansaction(ctx, session)
	}()
	err = p.productRepo.Save(ctx, product, session)
	if err != nil {
		return nil, err
	}
	err = p.transactionService.CommitTransaction(ctx, session)
	if err != nil {
		return nil, err
	}
	return &results.UpdateProductResult{Product: product}, nil
}

func (p *ProductCommandService) AddProductCategories(ctx context.Context, command *commands.AddProductCategoiesCommand) (*results.AddProductCategoriesResult, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, p.tracer, "Product.AddCategoreis: service", nil)
	defer tracing.EndSpan(span, err, nil)
	//find by id
	product, err := p.productRepo.FindById(ctx, command.ProductId)
	if err != nil {
		return nil, err
	}
	//add categories
	err = product.AddCategory(command.NewCategories)
	if err != nil {
		return nil, err
	}
	//save to db
	session, err := p.transactionService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			p.transactionService.RollbackTransaction(ctx, session)
		}
		p.transactionService.EndTransansaction(ctx, session)
	}()
	err = p.productRepo.Save(ctx, product, session)
	if err != nil {
		return nil, err
	}
	err = p.transactionService.CommitTransaction(ctx, session)
	if err != nil {
		return nil, err
	}
	return &results.AddProductCategoriesResult{Product: product}, nil
}

func (p *ProductCommandService) AddProductVariations(ctx context.Context, command *commands.AddProductVariationsCommand) (*results.AddProductVariationsResult, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, p.tracer, "Product.AddVariations: service", nil)
	defer tracing.EndSpan(span, err, nil)
	//find by id
	product, err := p.productRepo.FindById(ctx, command.ProductId)
	if err != nil {
		return nil, err
	}
	err = product.AddVariation(command.NewVariations)
	if err != nil {
		return nil, err
	}
	//save to db
	session, err := p.transactionService.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			p.transactionService.RollbackTransaction(ctx, session)
		}
		p.transactionService.EndTransansaction(ctx, session)
	}()
	err = p.productRepo.Save(ctx, product, session)
	if err != nil {
		return nil, err
	}
	err = p.transactionService.CommitTransaction(ctx, session)
	if err != nil {
		return nil, err
	}
	return &results.AddProductVariationsResult{Product: product}, nil
}

func NewProductCommandService(
	categoryRepo repositories.CategoryCommandRepository,
	productRepo repositories.ProductCommandRepository,
	shopService ShopService,
	mongoClient mongo.MongoTransactionService,
	tracer trace.Tracer,
	s3Config *config.S3Config,
	storageService storage.StorageService,
) handlers.ProductCommandHandler {
	return &ProductCommandService{
		categoryRepo:       categoryRepo,
		productRepo:        productRepo,
		shopService:        shopService,
		transactionService: mongoClient,
		tracer:             tracer,
		s3Config:           s3Config,
		storageService:     storageService,
	}
}
