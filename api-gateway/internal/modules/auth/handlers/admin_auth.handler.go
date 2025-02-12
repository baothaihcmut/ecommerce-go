package handlers

import (
	"go.opentelemetry.io/otel/trace"
)

type AdminAuthHandler struct {
	tracer trace.Tracer
}
