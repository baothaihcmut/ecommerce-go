package handlers

import (
	"context"
	"time"

	grpcpool "github.com/processout/grpc-go-pool"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/dtos/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/proto"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

type CategoryHandler interface {
	CreateCategory(ctx context.Context, request *request.CreateCategoryRequestDTO) (*response.CreateCategoryResponseDTO, error)
	BulkCreateCategories(ctx context.Context, dto *request.BulkCreateCategoriesRequestDTO) (_ *response.BulkCreateCategoriesResponseDTO, err error)
}

type CategoryHandlerImpl struct {
	productConnPool *grpcpool.Pool
	tracer          trace.Tracer
}

func NewCategoryHandler(productConnPool *grpcpool.Pool, tracer trace.Tracer) CategoryHandler {
	return &CategoryHandlerImpl{
		productConnPool: productConnPool,
		tracer:          tracer,
	}
}

func (c *CategoryHandlerImpl) CreateCategory(ctx context.Context, request *request.CreateCategoryRequestDTO) (_ *response.CreateCategoryResponseDTO, err error) {
	ctx, span := tracing.StartSpan(ctx, c.tracer, "Category.Create: Call product service", nil)
	defer tracing.EndSpan(span, err, nil)
	conn, err := c.productConnPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	s := proto.NewCategoryServiceClient(conn)
	req := proto.CreateCategoryRequest{
		Name:              request.Name,
		ParentCategoryIds: request.ParentCategoryIds,
	}
	res, err := s.CreateCategory(ctx, &req)
	if err != nil {
		return nil, err
	}
	return &response.CreateCategoryResponseDTO{
		Id:                res.Data.Id,
		Name:              res.Data.Name,
		ParentCategoryIds: res.Data.ParentCategoryIds,
	}, nil
}
func (c *CategoryHandlerImpl) BulkCreateCategories(ctx context.Context, dto *request.BulkCreateCategoriesRequestDTO) (_ *response.BulkCreateCategoriesResponseDTO, err error) {
	ctx, span := tracing.StartSpan(ctx, c.tracer, "Category.BulkCreate: Call product service", nil)
	defer tracing.EndSpan(span, err, nil)
	conn, err := c.productConnPool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	//map to proto
	categories := make([]*proto.CreateCategoryRequest, len(dto.Categories))
	for idx, category := range dto.Categories {
		categories[idx] = &proto.CreateCategoryRequest{
			Name:              category.Name,
			ParentCategoryIds: category.ParentCategoryIds,
		}
	}
	s := proto.NewCategoryServiceClient(conn)
	req := proto.BulkCreateCategoryRequest{
		Categories: categories,
	}
	//reponse time
	startTime := time.Now()

	res, err := s.BulkCreateCategory(ctx, &req)
	responseTime := time.Since(startTime)
	if err != nil {
		return nil, err
	}
	tracing.SetSpanAttribute(span, map[string]interface{}{
		"grpc_response_time": responseTime.String(),
	})

	result := make([]*response.CreateCategoryResponseDTO, len(res.Data))
	//map result
	for idx, category := range res.Data {
		result[idx] = &response.CreateCategoryResponseDTO{
			Id:                category.Id,
			Name:              category.Name,
			ParentCategoryIds: category.ParentCategoryIds,
		}
	}
	return &response.BulkCreateCategoriesResponseDTO{
		Categories: result,
	}, nil
}
