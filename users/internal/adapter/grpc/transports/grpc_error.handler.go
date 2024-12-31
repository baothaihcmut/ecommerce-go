package transports

import (
	"context"

	"github.com/go-kit/kit/transport"
	"github.com/go-kit/log"
)

type GrpcErrorHandler struct {
	logger log.Logger
}

func (h *GrpcErrorHandler) Handle(ctx context.Context, err error) {
	if err != nil {
		h.logger.Log("msg", "Request Failed", "err", err)
	}
}

func NewGrpcErrorHandler(logger log.Logger) transport.ErrorHandler {
	return &GrpcErrorHandler{
		logger: logger,
	}
}
