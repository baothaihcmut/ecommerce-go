package repositories

import (
	"context"
	"database/sql"
	"sync"

	"github.com/baothaihcmut/ecommerce-go/users/internal/adapter/db/postgres/sqlc"
	"github.com/baothaihcmut/ecommerce-go/users/internal/core/domain/entities"
	"github.com/baothaihcmut/ecommerce-go/users/internal/core/port/outbound/repositories"
	"github.com/google/uuid"
)

type PostgresUserRepo struct {
	q *sqlc.Queries
}

func toDBUser(user *entities.User) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		ID:                  uuid.NullUUID{UUID: user.Id, Valid: true},
		Email:               sql.NullString{String: user.Email, Valid: user.Email != ""},
		Password:            sql.NullString{String: user.Password, Valid: user.Password != ""},
		PhoneNumber:         sql.NullString{String: user.PhoneNumber, Valid: user.PhoneNumber != ""},
		IsShopOwnerActive:   sql.NullBool{Bool: user.IsShopOwnerActive, Valid: true},
		FirstName:           sql.NullString{String: user.FirstName, Valid: user.FirstName != ""},
		LastName:            sql.NullString{String: user.LastName, Valid: user.LastName != ""},
		CurrentRefreshToken: fromNullableString(user.CurrentRefreshToken),
	}
}
func fromNullableString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// CreateUser implements repositories.UserRepo.
func (p PostgresUserRepo) CreateUser(ctx context.Context, user *entities.User) error {
	//create in user table
	err := p.q.CreateUser(ctx, toDBUser(user))
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
			err = p.q.CreateCustomer(ctx, sqlc.CreateCustomerParams{
				UserId:     uuid.NullUUID{UUID: user.Id, Valid: true},
				LoyalPoint: sql.NullInt32{Int32: int32(user.Customer.LoyaltyPoint), Valid: true},
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
			err = p.q.CreateShopOwner(ctx, sqlc.CreateShopOwnerParams{
				UserId:           uuid.NullUUID{UUID: user.Id, Valid: true},
				BussinessLicense: sql.NullString{String: user.ShopOwner.BussinessLicense, Valid: user.ShopOwner.BussinessLicense != ""},
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
			err = p.q.CreateAddress(ctx, sqlc.CreateAddressParams{
				Priority: sql.NullInt32{Int32: int32(address.Priority), Valid: true},
				UserId:   uuid.NullUUID{UUID: user.Id, Valid: true},
				Street:   sql.NullString{String: address.Street, Valid: true},
				Town:     sql.NullString{String: address.Town, Valid: true},
				City:     sql.NullString{String: address.City, Valid: true},
				Province: sql.NullString{String: address.Province, Valid: true},
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

// FindUserByEmail implements repositories.UserRepo.
func (p PostgresUserRepo) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	res, err := p.q.FindUserByCriteria(ctx, sqlc.FindUserByCriteriaParams{
		Criteria: "phone_number",
		Value: sql.NullString{
			String: email,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, err
	}
	var currentRefreshToken *string
	if res.CurrentRefreshToken.Valid {
		currentRefreshToken = &res.CurrentRefreshToken.String
	}
	return &entities.User{
		Id:                  res.ID,
		Email:               res.Email,
		FirstName:           res.FirstName,
		LastName:            res.LastName,
		Password:            res.Password,
		PhoneNumber:         res.PhoneNumber,
		IsShopOwnerActive:   res.IsShopOwnerActive.Bool,
		CurrentRefreshToken: currentRefreshToken,
	}, nil
}

// FindUserByPhoneNumber implements repositories.UserRepo.
func (p PostgresUserRepo) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error) {
	res, err := p.q.FindUserByCriteria(ctx, sqlc.FindUserByCriteriaParams{
		Criteria: "phone_number",
		Value: sql.NullString{
			String: phoneNumber,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, err
	}
	var currentRefreshToken *string
	if res.CurrentRefreshToken.Valid {
		currentRefreshToken = &res.CurrentRefreshToken.String
	}
	return &entities.User{
		Id:                  res.ID,
		Email:               res.Email,
		FirstName:           res.FirstName,
		LastName:            res.LastName,
		Password:            res.Password,
		PhoneNumber:         res.PhoneNumber,
		IsShopOwnerActive:   res.IsShopOwnerActive.Bool,
		CurrentRefreshToken: currentRefreshToken,
	}, nil

}

func NewPostgresUserRepo(db *sql.DB) repositories.UserRepo {
	return PostgresUserRepo{
		q: sqlc.New(db),
	}
}
