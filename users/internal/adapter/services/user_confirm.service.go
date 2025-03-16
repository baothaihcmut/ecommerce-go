package services

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/external"
)

type UserConfirmService struct {
}

// SendEmail implements external.UserConfirmService.
func (u *UserConfirmService) SendEmail(context.Context, external.SendEmailArg) error {
	fmt.Println("Send email...")
	return nil
}

// StoreUserInfo implements external.UserConfirmService.
func (u *UserConfirmService) StoreUserInfo(context.Context, *entities.User) (string, error) {
	fmt.Println("Store user info...")
	return "", nil
}

func NewUserConfirmService() external.UserConfirmService {
	return &UserConfirmService{}
}
