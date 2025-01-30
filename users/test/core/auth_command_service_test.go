package core_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/models"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/commands"
	services "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/services/command"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJwtPort struct {
	mock.Mock
}

func (m *MockJwtPort) GenerateAccessToken(ctx context.Context, user *user.User) (valueobject.Token, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(valueobject.Token), args.Error(1)
}

func (m *MockJwtPort) GenerateRefreshToken(ctx context.Context, user *user.User) (valueobject.Token, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(valueobject.Token), args.Error(1)
}

func (m *MockJwtPort) DecodeAccessToken(ctx context.Context, token valueobject.Token) (models.AccessTokenSub, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(models.AccessTokenSub), args.Error(1)
}

func (m *MockJwtPort) DecodeRefreshToken(ctx context.Context, token valueobject.Token) (models.RefreshTokenSub, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(models.RefreshTokenSub), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

// FindById implements outbound.UserRepository.
func (m *MockUserRepository) FindById(ctx context.Context, id valueobject.UserId) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) CheckEmailExist(ctx context.Context, email valueobject.Email) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CheckPhoneNumberExist(ctx context.Context, phoneNumber valueobject.PhoneNumber) (bool, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, user *user.User, tx *sql.Tx) error {
	args := m.Called(ctx, user, tx)
	return args.Error(0)
}

func TestSignUp_Success(t *testing.T) {
	mockJwtPort := new(MockJwtPort)
	mockUserRepository := new(MockUserRepository)
	db := new(sql.DB)
	service := services.NewAuthCommandService(mockUserRepository, mockJwtPort, db)
	userDb := &user.User{
		Id:          valueobject.UserId(uuid.New()),
		Email:       valueobject.Email("baothai@gmail.com"),
		Password:    "baothai",
		PhoneNumber: "0828537679",
		Address: []valueobject.Address{
			{
				Priority: 0,
				Street:   "test",
				Town:     "test",
				City:     "test",
				Province: "test",
			},
		},
		Role:                enums.CUSTOMER,
		FirstName:           "thai",
		LastName:            "bao",
		CurrentRefreshToken: &valueobject.Token{Value: "refresh_token"},
		Customer: &entities.Customer{
			LoyaltyPoint: valueobject.LoyaltyPoint(0),
			Rank:         valueobject.Rank(enums.BRONZE),
		},
	}
	mockUserRepository.On("CheckEmailExist", mock.Anything, userDb.Email).Return(false, nil)
	mockUserRepository.On("CheckPhoneNumberExist", mock.Anything, userDb.PhoneNumber).Return(false, nil)
	mockUserRepository.On("Save", mock.Anything, userDb, mock.Anything).Return(nil)
	mockJwtPort.On("GenerateAccessToken", mock.Anything, mock.Anything).Return(valueobject.Token{Value: "access_token", TokenType: enums.ACCESS_TOKEN}, nil)
	mockJwtPort.On("GenerateRefreshToken", mock.Anything, mock.Anything).Return(valueobject.Token{Value: "refresh_token", TokenType: enums.REFRESH_TOKEN}, nil)

	command := &commands.SignUpCommand{
		Email:       string(userDb.Email),
		Password:    string(userDb.Password),
		PhoneNumber: string(userDb.PhoneNumber),
		Addresses: []*commands.Address{
			{
				Priority: userDb.Address[0].Priority,
				Street:   userDb.Address[0].Street,
				Town:     userDb.Address[0].Town,
				City:     userDb.Address[0].City,
				Province: userDb.Address[0].Province,
			},
		},
		FirstName:    userDb.FirstName,
		LastName:     userDb.LastName,
		Role:         userDb.Role,
		CustomerInfo: &commands.CustomerInfo{},
	}
	res, err := service.SignUp(context.Background(), command)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.AccessToken.Value, "access_token")
	assert.Equal(t, res.RefreshToken.Value, "refresh_token")
}
