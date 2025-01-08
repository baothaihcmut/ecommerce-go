package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/query/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/query/queries"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/query/results"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/outbound"
)

type UserQueryService struct {
	userRepo outbound.UserRepository
}

func (u *UserQueryService) FindUserById(ctx context.Context, q *queries.FindUserByIdQuery) (*results.UserQueryResult, error) {
	userId, err := valueobject.NewUserId(q.Id)
	if err != nil {
		return nil, err
	}
	doneCh := make(chan *user.User)
	errCh := make(chan error)
	go func() {
		user, err := u.userRepo.FindById(ctx, *userId)
		if err != nil {
			errCh <- err
		}
		doneCh <- user
	}()
	select {
	case user := <-doneCh:
		return &results.UserQueryResult{
			Id:          user.Id,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Customer:    user.Customer,
			ShopOwner:   user.ShopOwner,
		}, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func NewUserQueryService(userRepo outbound.UserRepository) handlers.UserQueryHandler {
	return &UserQueryService{
		userRepo: userRepo,
	}
}
