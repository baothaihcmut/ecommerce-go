package mappers

import (
	v1 "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/products/v1"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/inbound/results"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCreateProductCommand(req *v1.CreateProductRequest) *commands.CreateProductCommand {
	shopId, _ := primitive.ObjectIDFromHex(req.ShopId)
	return &commands.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		CategoryIds: lo.Map(req.CategoryIds, func(item string, _ int) primitive.ObjectID {
			categoryId, _ := primitive.ObjectIDFromHex(item)
			return categoryId
		}),
		ShopId:       shopId,
		NumOfImages:  int(req.NumOfImage),
		HasThumbNail: req.HasThumbnail,
		Variations:   req.Variations,
	}
}

func ToProductData(product *results.ProductResult) *v1.ProductData {
	return &v1.ProductData{
		Id:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		CategoryIds: lo.Map(product.CategoryIds, func(item primitive.ObjectID, _ int) string {
			return item.Hex()
		}),
		ShopId:     product.ShopId.Hex(),
		Variations: product.Variations,
		SoldTotal:  int32(product.SoldTotal),
		CreatedAt:  timestamppb.New(product.CreatedAt),
		UpdatedAt:  timestamppb.New(product.UpdatedAt),
	}
}
func ToUploadDetail(res *results.UploadImageResult) *v1.UploadDetail {
	return &v1.UploadDetail{
		Url:    res.Url,
		Method: res.Method,
		Expiry: int32(res.Expiry),
	}
}

func ToCreateProductReseponse(res *results.CreateProductResult) *v1.CreateProductData {
	return &v1.CreateProductData{
		Product:               ToProductData(&res.ProductResult),
		ThumbNailUploadDetail: ToUploadDetail(res.UploadThumbnailDetail),
		ImageUploadDetails: lo.Map(res.UploadImageDetails, func(item results.UploadImageResult, _ int) *v1.UploadDetail {
			return ToUploadDetail(&item)
		}),
	}
}
