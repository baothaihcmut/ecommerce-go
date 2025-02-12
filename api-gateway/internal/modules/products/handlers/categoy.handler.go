package handlers

import (
	"context"

	grpcpool "github.com/processout/grpc-go-pool"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/dtos/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/proto"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

type CategoryHandler interface {
	CreateCategory(ctx context.Context, request *request.CreateCategoryRequestDTO) (*response.CreateCategoryResponseDTO, error)
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

func (c *CategoryHandlerImpl) CreateCategory(ctx context.Context, request *request.CreateCategoryRequestDTO) (*response.CreateCategoryResponseDTO, error) {
	var err error
	ctx, span := tracing.StartSpan(ctx, c.tracer, "Category.Create: Call product service", nil)
	defer tracing.EndSpan(span, err, nil)
	conn, err := c.productConnPool.Get(ctx)
	defer conn.Close()
	if err != nil {
		return nil, err
	}
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
