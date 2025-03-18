package repositories

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/db/postgres/sqlc"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)



func toDBUser(user *entities.User) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		ID:                  pgtype.UUID{Bytes: user.Id, Valid: true},
		Email:               pgtype.Text{String: user.Email, Valid: user.Email != ""},
		Password:            pgtype.Text{String: user.Password, Valid: user.Password != ""},
		PhoneNumber:         pgtype.Text{String: user.PhoneNumber, Valid: user.PhoneNumber != ""},
		IsShopOwnerActive:   pgtype.Bool{Bool: user.IsShopOwnerActive, Valid: true},
		FirstName:           pgtype.Text{String: user.FirstName, Valid: user.FirstName != ""},
		LastName:            pgtype.Text{String: user.LastName, Valid: user.LastName != ""},
		CurrentRefreshToken: fromNullableString(user.CurrentRefreshToken),
	}
}
func fromNullableString(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// CreateUser implements repositories.UserRepo.
func (p PostgresUserRepo) CreateUser(ctx context.Context, user *entities.User, tx pgx.Tx) error {
	q := p.q
	if tx != nil {
		q = sqlc.New(tx)
	}
	//create in user table
	err := q.CreateUser(ctx, toDBUser(user))
	if err != nil {
		return err
	}
	//save sub enties
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if user.Customer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = q.CreateCustomer(ctx, sqlc.CreateCustomerParams{
				UserId:     pgtype.UUID{Bytes: user.Id, Valid: true},
				LoyalPoint: pgtype.Int4{Int32: int32(user.Customer.LoyaltyPoint), Valid: true},
			})
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancel()
				errCh <- err
			}
		}()
	}
	if user.ShopOwner != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = q.CreateShopOwner(ctx, sqlc.CreateShopOwnerParams{
				UserId:           pgtype.UUID{Bytes: user.Id, Valid: true},
				BussinessLicense: pgtype.Text{String: user.ShopOwner.BussinessLicense, Valid: user.ShopOwner.BussinessLicense != ""},
			})
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancel()
				errCh <- err
			}
		}()
	}
	for _, address := range user.Addresses {
		wg.Add(1)
		go func() {
			err = q.CreateAddress(ctx, sqlc.CreateAddressParams{
				Priority: pgtype.Int4{Int32: int32(address.Priority), Valid: true},
				UserId:   pgtype.UUID{Bytes: user.Id, Valid: true},
				Street:   pgtype.Text{String: address.Street, Valid: true},
				Town:     pgtype.Text{String: address.Town, Valid: true},
				City:     pgtype.Text{String: address.City, Valid: true},
				Province: pgtype.Text{String: address.Province, Valid: true},
			})
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancel()
				errCh <- err
			}
		}()
	}
	wg.Done()
	select {
	case err = <-errCh:
		return err
	default:
		return nil
	}
}