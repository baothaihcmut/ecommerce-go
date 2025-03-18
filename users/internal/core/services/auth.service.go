package services

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/db"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
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
	dbService 		   db.DBService
	logger 			   logger.Logger
}

// ConfirmSignUp implements handlers.AuthHandler.
func (a *AuthService) ConfirmSignUp(ctx context.Context, command *commands.ConfirmSignUpCommand) (*results.ConfirmSignUpResult, error) {
	//get user from cache
	user, err := a.userConfirmService.GetUserInfo(ctx, command.Code)
	if err != nil {
		return nil, err
	}
	//save user to db
	tx,err := a.dbService.BeginTransaction(ctx, db.DBTransactionReadWriteMode)
	if err != nil{
		return nil,err
	}
	defer func ()  {
		if err != nil{
			if err := a.dbService.RollBackTransaction(ctx,tx) ; err != nil {
				a.logger.Errorf(ctx,nil,"Error rollback transaction: %v", err)
			}
		}	
		if err := a.dbService.CommitTransaction(ctx,tx) ; err != nil {
			a.logger.Errorf(ctx,nil,"Error commit transaction: %v", err)
		}
	}()
	err = a.userRepo.CreateUser(ctx,user,tx)
	if err != nil{
		return nil,err
	}
	return &results.ConfirmSignUpResult{},nil
}

func (a *AuthService) SignUp(ctx context.Context, command *commands.SignUpCommand) (*results.SignUpResult, error) {
	//check if user pending for confirm
	isPending, err := a.userConfirmService.IsUserPendingConfirmSignUp(ctx, command.Email)
	if err != nil {
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
	logger logger.Logger,
) handlers.AuthHandler {
	return &AuthService{
		userRepo:           userRepo,
		jwtService:         jwtService,
		userConfirmService: userConfirmService,
		eventPublisher:     eventPublisher,
		logger: logger,
		dbService: dbService,
	}
}
