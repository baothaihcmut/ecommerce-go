package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/dtos/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/discovery"
)

type AuthHandler interface {
	LogIn(context.Context, request.LoginRequestDTO) (response.LoginResponeDTO, error)
}

type AuthHandlerImpl struct {
	discoveryService discovery.DiscoveryService
}

func (h *AuthHandlerImpl) LogIn(ctx context.Context, dto request.LoginRequestDTO) (response.LoginResponeDTO, error) {
	host, port, err := h.discoveryService.DiscoverService("user-service", "")
	if err != nil {
		return nil, err
	}

}
