package repositories

import (
	"context"
	"database/sql"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/sqlc/sqlc"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/outbound"
	"github.com/google/uuid"
)

type PostgresUserRepo struct {
	queries *sqlc.Queries
	conn    *sql.DB
}

// FindByEmail implements outbound.UserRepository.

func NewPostgresUserRepo(db *sql.DB) outbound.UserRepository {
	return &PostgresUserRepo{
		conn:    db,
		queries: sqlc.New(db),
	}
}

func (repo *PostgresUserRepo) toCreateUserArg(user *user.User) *sqlc.CreateUserParams {
	return &sqlc.CreateUserParams{
		ID:          uuid.NullUUID{UUID: uuid.UUID(user.Id), Valid: true},
		Email:       sql.NullString{String: string(user.Email), Valid: true},
		FirstName:   sql.NullString{String: user.FirstName, Valid: true},
		LastName:    sql.NullString{String: user.LastName, Valid: true},
		PhoneNumber: sql.NullString{String: string(user.PhoneNumber), Valid: true},
	}
}

func (repo *PostgresUserRepo) toCreateAddressArg(userId valueobject.UserId, address valueobject.Address) *sqlc.CreateAddressParams {
	return &sqlc.CreateAddressParams{
		Street:   sql.NullString{String: address.Street, Valid: true},
		Town:     sql.NullString{String: address.Town, Valid: true},
		City:     sql.NullString{String: address.City, Valid: true},
		Province: sql.NullString{String: address.Province, Valid: true},
		UserId:   uuid.NullUUID{UUID: uuid.UUID(userId), Valid: true},
	}
}

func (repo *PostgresUserRepo) toCreateCustomerArg(customer *entities.Customer) *sqlc.CreateCustomerParams {
	return &sqlc.CreateCustomerParams{
		UserId:     uuid.NullUUID{UUID: uuid.UUID(customer.Id), Valid: true},
		LoyalPoint: sql.NullInt32{Int32: int32(customer.LoyaltyPoint), Valid: true},
	}
}
func (repo *PostgresUserRepo) toCreateShopOwnerArg(shopOwner *entities.ShopOwner) *sqlc.CreateShopOwnerParams {
	return &sqlc.CreateShopOwnerParams{
		UserId:           uuid.NullUUID{UUID: uuid.UUID(shopOwner.Id), Valid: true},
		BussinessLicense: sql.NullString{String: shopOwner.BussinessLicense, Valid: true},
	}
}

func (repo *PostgresUserRepo) toUserDomain(result *sqlc.FindUserByIdRow, addresses []sqlc.Address) (*user.User, error) {
	userId, err := valueobject.NewUserId(result.ID)
	if err != nil {
		return nil, err
	}
	email, err := valueobject.NewEmail(result.Email)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := valueobject.NewPhoneNumber(result.PhoneNumber)
	if err != nil {
		return nil, err
	}
	userAddresses := make([]valueobject.Address, len(addresses))
	for _, address := range addresses {
		userAddresses = append(userAddresses, *valueobject.NewAddress(
			address.Street, address.Town, address.City, address.Province,
		))
	}
	var customer *entities.Customer
	var shopOwner *entities.ShopOwner
	if result.Role.RoleEnum == sqlc.RoleEnum(enums.CUSTOMER) {
		customer = entities.NewCustomerWithPoint(*userId, valueobject.LoyaltyPoint(result.LoyalPoint.Int32))
	} else {
		shopOwner = entities.NewShopOwner(*userId, result.BussinessLicense.String)
	}

	return &user.User{
		Id:          *userId,
		Email:       *email,
		PhoneNumber: *phoneNumber,
		Address:     userAddresses,
		Role:        enums.Role(result.Role.RoleEnum),
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Customer:    customer,
		ShopOwner:   shopOwner,
	}, nil
}

func (repo *PostgresUserRepo) Save(ctx context.Context, user *user.User, tx *sql.Tx) error {
	//create user
	err := repo.queries.WithTx(tx).CreateUser(ctx, *repo.toCreateUserArg(user))
	if err != nil {
		return err
	}
	//create in address table
	for _, address := range user.Address {
		err = repo.queries.WithTx(tx).CreateAddress(ctx, *repo.toCreateAddressArg(user.Id, address))
		if err != nil {
			return err
		}
	}

	//create in sub entity
	if user.Role == enums.CUSTOMER {
		err = repo.queries.WithTx(tx).CreateCustomer(ctx, *repo.toCreateCustomerArg(user.Customer))
		if err != nil {
			return err
		}
	} else {
		err = repo.queries.WithTx(tx).CreateShopOwner(ctx, *repo.toCreateShopOwnerArg(user.ShopOwner))
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *PostgresUserRepo) FindById(ctx context.Context, id valueobject.UserId) (*user.User, error) {
	userRes, err := repo.queries.FindUserById(ctx, uuid.NullUUID{UUID: uuid.UUID(id)})
	if err != nil {
		return nil, err
	}
	addressRes, err := repo.queries.FindAllAddressOfUser(ctx, uuid.NullUUID{UUID: uuid.UUID(id)})
	if err != nil {
		return nil, err
	}
	user, err := repo.toUserDomain(&userRes, addressRes)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (repo *PostgresUserRepo) FindByEmail(context.Context, valueobject.Email) (*user.User, error) {
	panic("unimplemented")
}

// FindByPhoneNumber implements outbound.UserRepository.
func (repo *PostgresUserRepo) FindByPhoneNumber(context.Context, valueobject.PhoneNumber) (*user.User, error) {
	panic("unimplemented")
}
