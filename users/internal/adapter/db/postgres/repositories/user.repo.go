package repositories

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/db/postgres/sqlc"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/repositories"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	q *sqlc.Queries
}

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
			err = p.q.CreateShopOwner(ctx, sqlc.CreateShopOwnerParams{
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
			err = p.q.CreateAddress(ctx, sqlc.CreateAddressParams{
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

// FindUserByPhoneNumber implements repositories.UserRepo.
func (p PostgresUserRepo) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error) {
	res, err := p.q.FindUserByCriteria(ctx, sqlc.FindUserByCriteriaParams{
		Criteria: "phone_number",
		Value: pgtype.Text{
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

func NewPostgresUserRepo(conn *pgxpool.Pool) repositories.UserRepo {
	return PostgresUserRepo{
		q: sqlc.New(conn),
	}
}
