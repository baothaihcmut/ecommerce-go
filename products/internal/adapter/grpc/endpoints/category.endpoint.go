package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	queryHandler "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/queries"
	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/trace"
)

type CategoryEndpoints struct {
	CreateCategory       endpoint.Endpoint
	FindAllCategory      endpoint.Endpoint
	BulkCreateCategories endpoint.Endpoint
}

func MakeCategoryEndpoints(c commandHandler.CategoryCommandHandler, q queryHandler.CategoryQueryHandler, tracer trace.Tracer) CategoryEndpoints {
	return CategoryEndpoints{
		CreateCategory:       makeCreateCategoryEndpoint(c, tracer),
		FindAllCategory:      makeFindAllCategoryEndpoint(q, tracer),
		BulkCreateCategories: makeBulkCreateCategoriesEndpoint(c, tracer),
	}
}

func makeCreateCategoryEndpoint(c commandHandler.CategoryCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		ctx, span := tracing.StartSpan(ctx, tracer, "Category.Create: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*commands.CreateCategoryCommand)
		res, err := c.CreateCategory(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeBulkCreateCategoriesEndpoint(c commandHandler.CategoryCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		ctx, span := tracing.StartSpan(ctx, tracer, "Category.BulkCreate: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*commands.BulkCreateCategories)
		res, err := c.BulkCreateCategories(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeFindAllCategoryEndpoint(q queryHandler.CategoryQueryHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		ctx, span := tracing.StartSpan(ctx, tracer, "Category.FindAll: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*queries.FindAllCategoryQuery)
		res, err := q.FindAllCategory(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
