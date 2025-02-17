package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/errors"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/sqlc/sqlc"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type PostgresUserRepo struct {
	queries *sqlc.Queries
	conn    *sql.DB
	tracer  trace.Tracer
}

// FindByEmail implements outbound.UserRepository.

func NewPostgresUserRepo(db *sql.DB, tracer trace.Tracer) outbound.UserRepository {
	return &PostgresUserRepo{
		conn:    db,
		queries: sqlc.New(db),
		tracer:  tracer,
	}
}

func (repo *PostgresUserRepo) toCreateUserArg(user *user.User) *sqlc.CreateUserParams {
	return &sqlc.CreateUserParams{
		ID:                  uuid.NullUUID{UUID: uuid.UUID(user.Id), Valid: true},
		Email:               sql.NullString{String: string(user.Email), Valid: true},
		Password:            sql.NullString{String: string(user.Password), Valid: true},
		FirstName:           sql.NullString{String: user.FirstName, Valid: true},
		LastName:            sql.NullString{String: user.LastName, Valid: true},
		PhoneNumber:         sql.NullString{String: string(user.PhoneNumber), Valid: true},
		CurrentRefreshToken: sql.NullString{String: user.CurrentRefreshToken.Value, Valid: (user.CurrentRefreshToken != nil)},
		Role:                sqlc.NullRoleEnum{RoleEnum: sqlc.RoleEnum(user.Role), Valid: true},
	}
}

func (repo *PostgresUserRepo) toCreateAddressArg(userId valueobject.UserId, address *entities.Address) *sqlc.CreateAddressParams {
	return &sqlc.CreateAddressParams{
		Priority: sql.NullInt32{Int32: int32(address.Priority), Valid: true},
		Street:   sql.NullString{String: address.Street, Valid: true},
		Town:     sql.NullString{String: address.Town, Valid: true},
		City:     sql.NullString{String: address.City, Valid: true},
		Province: sql.NullString{String: address.Province, Valid: true},
		UserId:   uuid.NullUUID{UUID: uuid.UUID(userId), Valid: true},
	}
}

func (repo *PostgresUserRepo) toCreateCustomerArg(userId uuid.UUID, customer *entities.Customer) *sqlc.CreateCustomerParams {
	return &sqlc.CreateCustomerParams{
		UserId:     uuid.NullUUID{UUID: userId, Valid: true},
		LoyalPoint: sql.NullInt32{Int32: int32(customer.LoyaltyPoint), Valid: true},
	}
}
func (repo *PostgresUserRepo) toCreateShopOwnerArg(userId uuid.UUID, shopOwner *entities.ShopOwner) *sqlc.CreateShopOwnerParams {
	return &sqlc.CreateShopOwnerParams{
		UserId:           uuid.NullUUID{UUID: userId, Valid: true},
		BussinessLicense: sql.NullString{String: shopOwner.BussinessLicense, Valid: true},
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
	var customer *entities.Customer
	var shopOwner *entities.ShopOwner
	if result.Role.RoleEnum == sqlc.RoleEnum(enums.CUSTOMER) {
		customer = entities.NewCustomerWithPoint(valueobject.LoyaltyPoint(result.LoyalPoint.Int32))
	} else {
		shopOwner = entities.NewShopOwner(result.BussinessLicense.String)
	}

	return &user.User{
		Id:          *userId,
		Email:       *email,
		Password:    valueobject.Password(result.Password),
		PhoneNumber: *phoneNumber,
		Address:     userAddresses,
		Role:        enums.Role(result.Role.RoleEnum),
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Customer:    customer,
		ShopOwner:   shopOwner,
	}, nil
}

func (repo *PostgresUserRepo) Save(ctx context.Context, user *user.User, tx *sql.Tx) (err error) {
	ctx, span := tracing.StartSpan(ctx, repo.tracer, "User.Save: database", nil)
	defer tracing.EndSpan(span, err, nil)
	//check if user exist
	_, err = repo.queries.CheckUserExistByCriteria(ctx, sqlc.CheckUserExistByCriteriaParams{
		Criteria: "id",
		Value:    sql.NullString{String: uuid.UUID(user.Id).String(), Valid: true},
	})
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var isUpdate bool
	if err == sql.ErrNoRows {
		isUpdate = false
	} else {
		isUpdate = true
	}

	//create user
	if isUpdate {
		err = repo.queries.WithTx(tx).UpdateUser(ctx, sqlc.UpdateUserParams{
			ID:          uuid.NullUUID{UUID: uuid.UUID(user.Id), Valid: true},
			Email:       sql.NullString{String: string(user.Email), Valid: true},
			Password:    sql.NullString{String: string(user.Password), Valid: true},
			FirstName:   sql.NullString{String: user.FirstName, Valid: true},
			LastName:    sql.NullString{String: user.LastName, Valid: true},
			PhoneNumber: sql.NullString{String: string(user.PhoneNumber), Valid: true},
			Role:        sqlc.NullRoleEnum{RoleEnum: sqlc.RoleEnum(user.Role), Valid: true},
		})
	} else {
		err = repo.queries.WithTx(tx).CreateUser(ctx, *repo.toCreateUserArg(user))
	}
	if err != nil {
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
			if isUpdate {
				err = repo.queries.WithTx(tx).UpdateAddress(ctx, sqlc.UpdateAddressParams{
					Street:   sql.NullString{String: address.Street, Valid: true},
					Town:     sql.NullString{String: address.Town, Valid: true},
					City:     sql.NullString{String: address.City, Valid: true},
					Province: sql.NullString{String: address.Province, Valid: true},
					UserId:   uuid.NullUUID{UUID: uuid.UUID(user.Id), Valid: true},
					Priority: sql.NullInt32{Int32: int32(address.Priority), Valid: true},
				})
			} else {
				err = repo.queries.WithTx(tx).CreateAddress(ctx, *repo.toCreateAddressArg(user.Id, address))
			}
			if err != nil {
				fmt.Println(err)
				if err == context.Canceled {
					return
				}
				cancel()
				errCh <- err
			}
		}()
	}

	if user.Role == enums.CUSTOMER {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if isUpdate {
				err = repo.queries.WithTx(tx).UpdateCustomer(ctx, sqlc.UpdateCustomerParams{
					LoyalPoint: sql.NullInt32{Int32: int32(user.Customer.LoyaltyPoint), Valid: true},
					UserId:     uuid.NullUUID{UUID: uuid.UUID(user.Id), Valid: true},
				})
			} else {
				err = repo.queries.WithTx(tx).CreateCustomer(ctx, *repo.toCreateCustomerArg(uuid.UUID(user.Id), user.Customer))
			}
			if err != nil {
				if err == context.Canceled {
					return
				}
				cancel()
				errCh <- err
			}
		}()
	}
	if user.Role == enums.SHOP_OWNER {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if isUpdate {
				err = repo.queries.WithTx(tx).UpdateShopOwner(ctx, sqlc.UpdateShopOwnerParams{
					BussinessLicense: sql.NullString{String: user.ShopOwner.BussinessLicense, Valid: true},
					UserId:           uuid.NullUUID{UUID: uuid.UUID(user.Id), Valid: true},
				})
			} else {
				err = repo.queries.WithTx(tx).CreateShopOwner(ctx, *repo.toCreateShopOwnerArg(uuid.UUID(user.Id), user.ShopOwner))
			}
			if err != nil {
				if err == context.Canceled {
					return
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
		fmt.Printf("Error:%v", err)
		return err
	default:
		return nil
	}

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
