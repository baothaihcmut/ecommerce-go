package services

import (
	"context"
	"sync"

	"github.com/baothaihcmut/ecommerce-go/users/internal/core/domain/entities"
	"github.com/baothaihcmut/ecommerce-go/users/internal/core/exception"
	"github.com/baothaihcmut/ecommerce-go/users/internal/core/port/inbound/commands"
	"github.com/baothaihcmut/ecommerce-go/users/internal/core/port/inbound/results"
	"github.com/baothaihcmut/ecommerce-go/users/internal/core/port/outbound/external"
	"github.com/baothaihcmut/ecommerce-go/users/internal/core/port/outbound/repositories"
	"github.com/samber/lo"
)

type AuthService struct {
	userRepo           repositories.UserRepo
	jwtService         external.JwtService
	userConfirmService external.UserConfirmService
}

func (a *AuthService) SignUp(ctx context.Context, command *commands.SignUpCommand) (*results.SignUpResult, error) {
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
	//send email
	err = a.userConfirmService.SendEmail(ctx, external.SendEmailArg{
		Email:     user.Email,
		Code:      code,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
	return &results.SignUpResult{}, nil

}
