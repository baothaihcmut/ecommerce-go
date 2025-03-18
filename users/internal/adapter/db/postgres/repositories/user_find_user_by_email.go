package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/db/postgres/sqlc"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// FindUserByEmail implements repositories.UserRepo.
func (p PostgresUserRepo) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	res, err := p.q.FindUserByCriteria(ctx, sqlc.FindUserByCriteriaParams{
		Criteria: "phone_number",
		Value: pgtype.Text{
			String: email,
			Valid:  true,
		},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	var currentRefreshToken *string
	if res.CurrentRefreshToken.Valid {
		currentRefreshToken = &res.CurrentRefreshToken.String
	}
	return &entities.User{
		Id:                  res.ID.Bytes,
		Email:               res.Email,
		FirstName:           res.FirstName,
		LastName:            res.LastName,
		Password:            res.Password,
		PhoneNumber:         res.PhoneNumber,
		IsShopOwnerActive:   res.IsShopOwnerActive,
		CurrentRefreshToken: currentRefreshToken,
	}, nil
}
