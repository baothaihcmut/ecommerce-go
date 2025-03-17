package services

import (
	"context"
	"fmt"
	"time"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/cache"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/external"
	"github.com/google/uuid"
)

type UserConfirmService struct {
	cacheService cache.CacheService
}

// GetUserInfo implements external.UserConfirmService.
func (u *UserConfirmService) GetUserInfo(ctx context.Context, code string) (*entities.User, error) {
	var user entities.User
	if err := u.cacheService.GetValue(ctx, fmt.Sprintf("user_sign_up_info:%s", code), &user); err != nil {
		return nil, err
	}
	//remove in cache
	if err := u.cacheService.Remove(ctx, fmt.Sprintf("user_sign_up_info:%s", code)); err != nil {
		return nil, err
	}
	if err := u.cacheService.Remove(ctx, fmt.Sprintf("user_pending_confirm_email:%s", user.Email)); err != nil {
		return nil, err
	}
	return &user, nil
}

// IsUserPendingConfirmSignUp implements external.UserConfirmService.
func (u *UserConfirmService) IsUserPendingConfirmSignUp(ctx context.Context, email string) (bool, error) {
	res, err := u.cacheService.GetString(ctx, fmt.Sprintf("user_pending_confirm_email:%s", email))
	if err != nil {
		return false, err
	}
	return res == "1", nil
}

// GenerateUrlForConfirm implements external.UserConfirmService.
func (u *UserConfirmService) GenerateUrlForConfirm(_ context.Context, code string) (string, error) {
	return fmt.Sprintf("http://localhost:3000/buyer?code=%s", code), nil
}

// StoreUserInfo implements external.UserConfirmService.
func (u *UserConfirmService) StoreUserInfo(ctx context.Context, user *entities.User) (string, error) {
	//store user email
	err := u.cacheService.SetString(ctx, fmt.Sprintf("user_pending_confirm_email:%s", user.Email), "1", time.Minute*30)
	if err != nil {
		return "", err
	}
	//generate
	code := uuid.New().String()
	//store user info
	err = u.cacheService.SetValue(ctx, fmt.Sprintf("user_sign_up_info:%s", code), user, time.Minute*30)
	if err != nil {
		return "", err
	}
	return code, nil
}

func NewUserConfirmService(
	cacheService cache.CacheService,
) external.UserConfirmService {
	return &UserConfirmService{
		cacheService: cacheService,
	}
}
