package middleware

import (
	"runtime"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

func RecoverMiddleware(tracer trace.Tracer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ctx, span := tracing.StartSpan(c.Request().Context(), tracer, "Api gateway: recover middleware", nil)
			defer tracing.EndSpan(span, err, nil)
			defer func() {
				if r := recover(); r != nil {
					switch r.(type) {
					case string:
						tracing.SetSpanAttribute(span, map[string]interface{}{
							"detail": r.(string),
						})
					case error:
						err = r.(error)
						tracing.SetSpanAttribute(span, map[string]interface{}{
							"detail": r.(error).Error(),
						})
					}
					buf := make([]byte, 1024)
					n := runtime.Stack(buf, false)
					// tracing
					tracing.SetSpanAttribute(span, map[string]interface{}{
						"level":      "FATAL",
						"stackTrace": string(buf[:n]),
					})
				}
			}()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
