package services

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/db"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/entities"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/events"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/exception"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/external"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/repositories"
	"github.com/samber/lo"
)

type AuthService struct {
	userRepo           repositories.UserRepo
	jwtService         external.JwtService
	userConfirmService external.UserConfirmService
	eventPublisher     external.EventPublisherService
	dbService          db.DBService
	queueService       queue.QueueService
	logger             logger.Logger
}

// LogIn implements handlers.AuthHandler.
func (a *AuthService) LogIn(ctx context.Context,command *commands.LogInCommand) (*results.LogInResult, error) {
	//find user by email
	user,err := a.userRepo.FindUserByEmail(ctx,command.Email)
	if err!= nil {
		a.logger.WithCtx(ctx).Errorf(map[string]any{
			"email":command.Email,
		},"Error find user by email: ",err)
	}
	if user == nil {
		return nil, exception.ErrWrongEmailOrPassword
	}
	if !user.ValidatePassword(command.Password) {
		return nil, exception.ErrWrongEmailOrPassword
	}
	//gen access and refresh token
	accessToken,err := a.jwtService.GenerateAccessToken(ctx, external.AccessTokenPayload{
		UserId: user.Id,
		IsShopOwnerActive: user.IsShopOwnerActive,
	})
	if err != nil{
		a.logger.WithCtx(ctx).Errorf(map[string]any{
			"email":command.Email,
		},"Error generate access token: ",err)
		return nil,err
	}
	refreshToken,err := a.jwtService.GenerateRefreshToken(ctx, external.RefreshTokenPayload{
		UserId: user.Id,
	})
	if err != nil{
		a.logger.WithCtx(ctx).Errorf(map[string]any{
			"email":command.Email,
		},"Error generate refesh token: ",err)
		return nil,err
	}
	return &results.LogInResult{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	},nil
}

// ConfirmSignUp implements handlers.AuthHandler.
func (a *AuthService) ConfirmSignUp(ctx context.Context, command *commands.ConfirmSignUpCommand) (*results.ConfirmSignUpResult, error) {
	//get user from cache
	user, err := a.userConfirmService.GetUserInfo(ctx, command.Code)
	if err != nil {
		return nil, err
	}
	//save user to db
	tx, err := a.dbService.BeginTransaction(ctx, db.DBTransactionReadWriteMode)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if err := a.dbService.RollBackTransaction(ctx, tx); err != nil {
				a.logger.WithCtx(ctx).Errorf( nil, "Error rollback transaction: %v", err)
			}
		} else {
			if err := a.dbService.CommitTransaction(ctx, tx); err != nil {
				a.logger.WithCtx(ctx).Errorf( nil, "Error commit transaction: %v", err)
			}
		}
	}()
	err = a.userRepo.WithTx(tx).CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &results.ConfirmSignUpResult{}, nil
}

func (a *AuthService) SignUp(ctx context.Context, command *commands.SignUpCommand) (*results.SignUpResult, error) {
	//check if user pending for confirm
	isPending, err := a.userConfirmService.IsUserPendingConfirmSignUp(ctx, command.Email)
	if err != nil {
		a.logger.WithCtx(ctx).Errorf( map[string]any{
			"email": command.Email,
		}, "Error check user pending sign up: %v", err)
		return nil, err
	}
	if isPending {
		return nil, exception.ErrUserPendingForConfirm
	}

	//check if user phone number or email exist
	wgCheck := sync.WaitGroup{}
	errCheck := make(chan error, 1)
	ctx, cancelCheck := context.WithCancel(ctx)
	defer cancelCheck()
	wgCheck.Add(1)
	go func() {
		defer wgCheck.Done()
		user, err := a.userRepo.FindUserByEmail(ctx, command.Email)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			a.logger.WithCtx(ctx).Errorf( nil, "Error get user by email: %v", err)
			cancelCheck()
			errCheck <- err
		}
		if user != nil {
			errCheck <- exception.ErrEmailExist
		}
	}()
	wgCheck.Add(1)
	go func() {
		defer wgCheck.Done()
		user, err := a.userRepo.FindUserByPhoneNumber(ctx, command.PhoneNumber)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			a.logger.WithCtx(ctx).Errorf( nil, "Error get user by phone number: %v", err)
			cancelCheck()
			errCheck <- err
		}
		if user != nil {
			errCheck <- exception.ErrPhonenumberExist
		}
	}()
	wgCheck.Wait()
	select {
	case err := <-errCheck:
		return nil, err
	default:
	}
	close(errCheck)
	user := entities.NewUser(
		command.Email,
		command.Password,
		command.PhoneNumber,
		lo.Map(command.Addresses, func(item commands.SignUpAddress, _ int) entities.CreateAddessArg {
			return entities.CreateAddessArg{
				Priority: item.Priority,
				Town:     item.Town,
				City:     item.City,
				Street:   item.Street,
				Province: item.Province,
			}
		}),
		command.FirstName,
		command.LastName,
	)
	//store user info
	code, err := a.userConfirmService.StoreUserInfo(ctx, user)
	if err != nil {
		a.logger.WithCtx(ctx).Errorf( nil, "Error store user to cache: %v", err)

		return nil, err
	}
	//generate url for confrim
	url, err := a.userConfirmService.GenerateUrlForConfirm(ctx, code)
	if err != nil {
		return nil, err
	}
	//publish event
	e := events.UserSignUpEvent{
		User:       user,
		ConfrimUrl: url,
	}
	if err := a.eventPublisher.PublishUserSignUpEvent(ctx, e); err != nil {
		a.logger.WithCtx(ctx).Errorf( nil, "Error publish event: %v", err)
		return nil, err
	}
	return &results.SignUpResult{}, nil

}

func NewAuthService(
	userRepo repositories.UserRepo,
	jwtService external.JwtService,
	userConfirmService external.UserConfirmService,
	eventPublisher external.EventPublisherService,
	dbService db.DBService,
	queueService queue.QueueService,
	logger logger.Logger,
) handlers.AuthHandler {
	return &AuthService{
		userRepo:           userRepo,
		jwtService:         jwtService,
		userConfirmService: userConfirmService,
		eventPublisher:     eventPublisher,
		logger:             logger,
		queueService:       queueService,
		dbService:          dbService,
	}
}

func (a *AuthService) RunBootstrap() {
	//init exchange
	a.queueService.InitExchange("user-events", "topic")
	a.queueService.BindQueue("mail-user-signup", "user.signup", "user-events")
}
