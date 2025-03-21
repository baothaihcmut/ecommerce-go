package services

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/constant"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/models"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/mongo"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/storage"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/domain/entities"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/exception"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/outbound/external"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/outbound/repositories"
)

type ProductService struct {
	productRepo    repositories.ProductRepo
	categoryRepo   repositories.CategoryRepo
	shopService    external.ShopService
	logger         logger.Logger
	mongoService   mongo.MongoTransactionService
	storageService storage.StorageService
}

func NewProductService(
	productRepo repositories.ProductRepo,
	categoryRepo repositories.CategoryRepo,
	shopService external.ShopService,
	logger logger.Logger,
	mongoService mongo.MongoTransactionService,
	storageService storage.StorageService,
) *ProductService {
	return &ProductService{
		productRepo:    productRepo,
		categoryRepo:   categoryRepo,
		shopService:    shopService,
		logger:         logger,
		mongoService:   mongoService,
		storageService: storageService,
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, command *commands.CreateProductCommand) (*results.CreateProductResult, error) {
	//get user Context
	userCtx := ctx.Value(constant.UserContextKey).(*models.UserContext)
	if !userCtx.IsShopOwnerActive {
		return nil, exception.ErrUserNotShopOwnerActive
	}
	//find shop
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		shop, err := p.shopService.FindShopById(ctx, command.ShopId)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			p.logger.WithCtx(ctx).Errorf(map[string]any{
				"shop_id": command.ShopId,
			}, "Error get shop by id: ", err)
			cancel()
			errCh <- err
			return

		}
		//shop not exist
		if shop == nil {
			errCh <- exception.ErrShopNotExist
			return
		}
		//user is not owner of shop
		if shop.ShopOwnerId != userCtx.UserId {
			errCh <- exception.ErrUserIsNotShopOwner
			return
		}
	}()
	//check category
	for _, categoryId := range command.CategoryIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			category, err := p.categoryRepo.FindCategoryById(ctx, categoryId)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				p.logger.WithCtx(ctx).Errorf(map[string]any{
					"shop_id": command.ShopId,
				}, "Error find category by id: ", err)
				cancel()
				errCh <- err
				return
			}
			if category == nil {
				errCh <- exception.ErrCategoryNotExist
				return
			}
		}()
	}
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, err
	default:
	}
	//create new product
	product := entities.NewProduct(
		command.Name,
		command.Description,
		command.ShopId,
		command.CategoryIds,
		command.Variations,
		command.HasThumbNail,
		command.NumOfImages,
	)
	//save to db
	session, err := p.mongoService.BeginTransaction(ctx)
	if err != nil {
		p.logger.WithCtx(ctx).Errorf(nil, "Error start mongo transaction: ", err)
		return nil, err
	}
	defer p.mongoService.EndTransansaction(ctx, session)
	if err := p.productRepo.CreateProduct(ctx, product); err != nil {
		p.logger.WithCtx(ctx).Errorf(nil, "Error create product: ", err)
		if err := p.mongoService.RollbackTransaction(ctx, session); err != nil {
			p.logger.WithCtx(ctx).Errorf(nil, "Error rollback mongo transaction: ", err)
		}
		return nil, err
	}
	//commit transaction
	if err := p.mongoService.CommitTransaction(ctx, session); err != nil {
		p.logger.WithCtx(ctx).Errorf(nil, "Error rollback mongo transaction: ", err)
		return nil, err
	}
	//get put image url
	wgStorage := sync.WaitGroup{}
	imageUrlCh := make(chan results.UploadImageResult, command.NumOfImages)
	uploadDetail := make([]results.UploadImageResult, 0, command.NumOfImages)
	errChStorage := make(chan error, 1)
	var thumbnailUploadDetail *results.UploadImageResult
	if product.Thumbnail != nil {
		wgStorage.Add(1)
		go func() {
			defer wgStorage.Done()
			url, err := p.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArg{
				Method: storage.GetPresignUrlMethodPut,
				Key:    *product.Thumbnail,
			})
			if err != nil {
				errChStorage <- err
			}
			thumbnailUploadDetail = &results.UploadImageResult{
				Url:    url,
				Method: "PUT",
				Expiry: 30,
			}
		}()
	}
	for i := 0; i < command.NumOfImages; i++ {
		wgStorage.Add(1)
		go func() {
			defer wgStorage.Done()
			url, err := p.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArg{
				Method: storage.GetPresignUrlMethodPut,
				Key:    *product.Thumbnail,
			})
			if err != nil {
				errChStorage <- err
			}
			imageUrlCh <- results.UploadImageResult{
				Url:    url,
				Method: "PUT",
				Expiry: 30,
			}
		}()
	}
	var errStorage error
	go func() {
		wgStorage.Wait()
		close(imageUrlCh)
	}()
	go func() {
		errStorage = <-errChStorage
		close(imageUrlCh)
	}()
	for detail := range imageUrlCh {
		uploadDetail = append(uploadDetail, detail)
	}
	if errStorage != nil {
		return nil, err
	}
	return &results.CreateProductResult{
		ProductResult:         results.MapToProductResult(product),
		UploadThumbnailDetail: thumbnailUploadDetail,
		UploadImageDetails:    uploadDetail,
	}, nil

}
