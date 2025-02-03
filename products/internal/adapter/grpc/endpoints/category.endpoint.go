package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/grpc/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/handlers"
	queryHandler "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/queries"
	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
)

type CategoryEndpoints struct {
	CreateCategory       endpoint.Endpoint
	FindAllCategory      endpoint.Endpoint
	BulkCreateCategories endpoint.Endpoint
}

func MakeCategoryEndpoints(c commandHandler.CategoryCommandHandler, q queryHandler.CategoryQueryHandler, tracer opentracing.Tracer) CategoryEndpoints {
	return CategoryEndpoints{
		CreateCategory:       middlewares.ExtractTracingMiddleware(tracer, "Create category")(makeCreateCategoryEndpoint(c)),
		FindAllCategory:      middlewares.ExtractTracingMiddleware(tracer, "Find all categories")(makeFindAllCategoryEndpoint(q)),
		BulkCreateCategories: middlewares.ExtractTracingMiddleware(tracer, "Bulk create categories")(makeBulkCreateCategoriesEndpoint(c)),
	}
}

func makeCreateCategoryEndpoint(c commandHandler.CategoryCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.CreateCategoryCommand)
		res, err := c.CreateCategory(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeBulkCreateCategoriesEndpoint(c commandHandler.CategoryCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.BulkCreateCategories)
		res, err := c.BulkCreateCategories(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeFindAllCategoryEndpoint(q queryHandler.CategoryQueryHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*queries.FindAllCategoryQuery)
		res, err := q.FindAllCategory(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
