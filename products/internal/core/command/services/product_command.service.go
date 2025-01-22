package services

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/exceptions"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductCommandService struct {
	productRepo  repositories.ProductCommandRepository
	categoryRepo repositories.CategoryCommandRepository
	shopService  ShopService
	mongo        *mongo.Client
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
	//check contraint
	err := p.checkContraints(ctx, product)
	if err != nil {
		return nil, err
	}
	productDomain, err := products.NewProduct(
		product.Name,
		product.Description,
		product.Unit,
		product.ShopId,
		product.CategoryIds,
		product.Variations,
	)
	if err != nil {
		return nil, err
	}
	//save to db
	session, err := p.mongo.StartSession()
	if err != nil {
		return nil, err
	}
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			session.AbortTransaction(ctx)
		}
	}()
	err = p.productRepo.Save(ctx, productDomain, session)
	if err != nil {
		return nil, err
	}
	session.CommitTransaction(ctx)
	return &results.CreateProductResult{
		Product: productDomain,
	}, nil
}
func NewProductCommandService(
	categoryRepo repositories.CategoryCommandRepository,
	productRepo repositories.ProductCommandRepository,
	shopService ShopService,
	mongoClient *mongo.Client,
) handlers.ProductCommandHandler {
	return &ProductCommandService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
		shopService:  shopService,
		mongo:        mongoClient,
	}
}
