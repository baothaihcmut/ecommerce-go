package repositories

import (
	"context"
	"database/sql"
	"sync"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/errors"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/sqlc/sqlc"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound/repositories"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type PostgresUserRepo struct {
	queries *sqlc.Queries
	conn    *sql.DB
	tracer  trace.Tracer
}

func NewPostgresUserRepo(db *sql.DB, tracer trace.Tracer) repositories.UserRepository {
	return &PostgresUserRepo{
		conn:    db,
		queries: sqlc.New(db),
		tracer:  tracer,
	}
}

func (repo *PostgresUserRepo) toUserDomain(result *sqlc.FindUserByCriteriaRow, addresses []sqlc.Address) (*user.User, error) {
	userId, err := valueobject.NewUserId(result.ID)
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	email, err := valueobject.NewEmail(result.Email)
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	phoneNumber, err := valueobject.NewPhoneNumber(result.PhoneNumber)
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	userAddresses := make([]*entities.Address, len(addresses))
	for _, address := range addresses {
		userAddresses = append(userAddresses, entities.NewAddress(
			int(address.Priority), address.Street, address.Town, address.City, address.Province,
		))
	}

	var shopOwner *entities.ShopOwner
	if result.IsShopOwnerActive.Bool {
		shopOwner = entities.NewShopOwner(result.BussinessLicense.String)
	}

	return &user.User{
		Id:                *userId,
		Email:             *email,
		Password:          valueobject.Password(result.Password),
		PhoneNumber:       *phoneNumber,
		Address:           userAddresses,
		IsShopOwnerActive: false,
		FirstName:         result.FirstName,
		LastName:          result.LastName,
		Customer:          entities.NewCustomerWithPoint(valueobject.LoyaltyPoint(result.LoyalPoint.Int32)),
		ShopOwner:         shopOwner,
	}, nil
}

func (repo *PostgresUserRepo) UpdateUserAddress(ctx context.Context, user *user.User, t *sql.Tx) (err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.Save: database", nil)
	defer tracing.EndSpan(span, err, nil)

	// find all address of user
	address, err := repo.queries.FindUserAddress(ctx, uuid.NullUUID{
		UUID:  uuid.UUID(user.Id),
		Valid: true,
	})
	if err != nil {
		return err
	}
	//map for check address with priority exist
	mapAddr := make(map[int]struct{})
	for _, addr := range address {
		mapAddr[int(addr.Priority)] = struct{}{}
	}
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 1)
	for _, address := range user.Address {
		//check if address exist
		wg.Add(1)
		go func() {
			//check if address exist
			if _, exist := mapAddr[address.Priority]; exist {
				//if exist update old
				err = repo.queries.UpdateAddress(ctx, sqlc.UpdateAddressParams{
					UserId: uuid.NullUUID{
						UUID:  uuid.UUID(user.Id),
						Valid: true,
					},
					Priority: sql.NullInt32{Int32: int32(address.Priority), Valid: true},
					Street:   sql.NullString{String: address.Street, Valid: true},
					Town:     sql.NullString{String: address.Town, Valid: true},
					City:     sql.NullString{String: address.City, Valid: true},
					Province: sql.NullString{String: address.Province, Valid: true},
				})
			} else {
				err = repo.queries.CreateAddress(ctx, sqlc.CreateAddressParams{
					UserId: uuid.NullUUID{
						UUID:  uuid.UUID(user.Id),
						Valid: true,
					},
					Priority: sql.NullInt32{Int32: int32(address.Priority), Valid: true},
					Street:   sql.NullString{String: address.Street, Valid: true},
					Town:     sql.NullString{String: address.Town, Valid: true},
					City:     sql.NullString{String: address.City, Valid: true},
					Province: sql.NullString{String: address.Province, Valid: true},
				})
			}
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
	wg.Wait()
	select {
	case err = <-errCh:
		return err
	default:
		return nil
	}
}

func (repo *PostgresUserRepo) CreateUser(ctx context.Context, user *user.User, tx *sql.Tx) (err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.Save: database", nil)
	defer tracing.EndSpan(span, err, nil)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//create user

	err = repo.queries.WithTx(tx).CreateUser(ctx, sqlc.CreateUserParams{
		ID:                  uuid.NullUUID{UUID: uuid.UUID(user.Id), Valid: true},
		Email:               sql.NullString{String: string(user.Email), Valid: true},
		Password:            sql.NullString{String: string(user.Password), Valid: true},
		FirstName:           sql.NullString{String: user.FirstName, Valid: true},
		LastName:            sql.NullString{String: user.LastName, Valid: true},
		PhoneNumber:         sql.NullString{String: string(user.PhoneNumber), Valid: true},
		CurrentRefreshToken: sql.NullString{String: user.CurrentRefreshToken, Valid: true},
	})
	if err != nil {
		tracing.SetSpanAttribute(span, map[string]interface{}{
			"detail": "insert user",
		})
		return err
	}
	//insert on sub table
	wg := sync.WaitGroup{}
	errCh := make(chan error, 1)
	//create in address table
	for _, address := range user.Address {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = repo.queries.WithTx(tx).CreateAddress(ctx, sqlc.CreateAddressParams{
				Priority: sql.NullInt32{Int32: int32(address.Priority), Valid: true},
				Street:   sql.NullString{String: address.Street, Valid: true},
				Town:     sql.NullString{String: address.Town, Valid: true},
				City:     sql.NullString{String: address.City, Valid: true},
				Province: sql.NullString{String: address.Province, Valid: true},
				UserId:   uuid.NullUUID{UUID: uuid.UUID(user.Id), Valid: true},
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

	//insert into customer table
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = repo.queries.WithTx(tx).CreateCustomer(ctx, sqlc.CreateCustomerParams{
			UserId: uuid.NullUUID{
				UUID:  uuid.UUID(user.Id),
				Valid: true,
			},
			LoyalPoint: sql.NullInt32{
				Int32: int32(user.Customer.LoyaltyPoint),
				Valid: true,
			},
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

	//if shopowner active insert into shop owner table
	if user.IsShopOwnerActive {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = repo.queries.WithTx(tx).CreateShopOwner(ctx, sqlc.CreateShopOwnerParams{
				UserId: uuid.NullUUID{
					UUID:  uuid.UUID(user.Id),
					Valid: true,
				},
				BussinessLicense: sql.NullString{
					String: user.ShopOwner.BussinessLicense,
					Valid:  true,
				},
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
	//wait for all worker
	wg.Wait()
	//if have error return error
	select {
	case err = <-errCh:
		return err
	default:
		return nil
	}
}

func (repo *PostgresUserRepo) ActivateShopOwner(ctx context.Context, user *user.User, t *sql.Tx) (err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.ActivateShopOwner: database", nil)
	defer tracing.EndSpan(span, err, nil)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 1)

	//update in user table
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = repo.queries.WithTx(t).UpdateUser(ctx, sqlc.UpdateUserParams{
			IsShopOwnerActive: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
			ID: uuid.NullUUID{
				UUID:  uuid.UUID(user.Id),
				Valid: true,
			},
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
	//create new in shop owner table
	wg.Add(1)
	go func() {
		err = repo.queries.WithTx(t).CreateShopOwner(
			ctx,
			sqlc.CreateShopOwnerParams{
				UserId: uuid.NullUUID{
					UUID:  uuid.UUID(user.Id),
					Valid: true,
				},
			},
		)
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
	wg.Wait()
	select {
	case err = <-errCh:
		return err
	default:
		return nil
	}
}

func (repo *PostgresUserRepo) UpdateUser(ctx context.Context, user *user.User, t *sql.Tx) (err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.UpdateInfo: database", nil)
	defer tracing.EndSpan(span, err, nil)
	err = repo.queries.WithTx(t).UpdateUser(ctx, sqlc.UpdateUserParams{
		Email: sql.NullString{
			String: string(user.Email),
			Valid:  true,
		},
		PhoneNumber: sql.NullString{
			String: string(user.PhoneNumber),
			Valid:  true,
		},
		FirstName: sql.NullString{
			String: user.FirstName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: user.LastName,
			Valid:  true,
		},
		ID: uuid.NullUUID{
			Valid: true,
			UUID:  uuid.UUID(user.Id),
		},
		CurrentRefreshToken: sql.NullString{
			Valid:  true,
			String: user.CurrentRefreshToken,
		},
		Password: sql.NullString{
			Valid:  true,
			String: string(user.Password),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// FindByEmail implements outbound.UserRepository.
func (repo *PostgresUserRepo) FindByEmail(ctx context.Context, email valueobject.Email) (resp *user.User, err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.FindByEmail: database", nil)
	defer tracing.EndSpan(span, err, nil)
	userRes, err := repo.queries.FindUserByCriteria(ctx, sqlc.FindUserByCriteriaParams{
		Criteria: "email",
		Value: sql.NullString{
			String: string(email),
			Valid:  true,
		},
	})
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	addressRes, err := repo.queries.FindAllAddressOfUser(ctx, uuid.NullUUID{UUID: userRes.ID})
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	user, err := repo.toUserDomain(&userRes, addressRes)
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	return user, nil
}
func (repo *PostgresUserRepo) FindById(ctx context.Context, id valueobject.UserId) (resp *user.User, err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.FindById: database", nil)
	defer tracing.EndSpan(span, err, nil)
	userRes, err := repo.queries.FindUserByCriteria(ctx, sqlc.FindUserByCriteriaParams{
		Criteria: "id",
		Value: sql.NullString{
			String: uuid.UUID(id).String(),
			Valid:  true,
		},
	})
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	addressRes, err := repo.queries.FindAllAddressOfUser(ctx, uuid.NullUUID{UUID: uuid.UUID(id)})
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	user, err := repo.toUserDomain(&userRes, addressRes)
	if err != nil {
		return nil, errors.NewError(err, errors.CaptureStackTrace())
	}
	return user, nil

}

func (repo *PostgresUserRepo) CheckEmailExist(ctx context.Context, email valueobject.Email) (resp bool, err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.CheckEmailExist: database", nil)
	defer tracing.EndSpan(span, err, nil)
	_, err = repo.queries.CheckUserExistByCriteria(ctx, sqlc.CheckUserExistByCriteriaParams{
		Criteria: "email",
		Value:    sql.NullString{String: string(email), Valid: true},
	})

	if err != nil {
		//if no row return return false
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *PostgresUserRepo) CheckPhoneNumberExist(ctx context.Context, phoneNumber valueobject.PhoneNumber) (resp bool, err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.CheckPhonenumberExist: database", nil)
	defer tracing.EndSpan(span, err, nil)
	_, err = repo.queries.CheckUserExistByCriteria(ctx, sqlc.CheckUserExistByCriteriaParams{
		Criteria: "phone_number",
		Value:    sql.NullString{String: string(phoneNumber), Valid: true},
	})
	if err != nil {
		//if no row return return false
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
