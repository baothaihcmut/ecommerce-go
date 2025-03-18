package mappers

import (
	v1 "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/results"
)

func ToSignUpCommand(proto *v1.SignUpRequest) *commands.SignUpCommand {
	return &commands.SignUpCommand{
		Email:       proto.Email,
		Password:    proto.Password,
		FirstName:   proto.FirstName,
		LastName:    proto.LastName,
		PhoneNumber: proto.PhoneNumber,
	}
}

func ToSignUpResponse(result *results.SignUpResult) *v1.SignUpData {
	return &v1.SignUpData{}
}

func ToConfirmSignUpCommand(proto *v1.ConfirmSignUpRequest) *commands.ConfirmSignUpCommand {
	return &commands.ConfirmSignUpCommand{
		Code: proto.Code,
	}
}
func ToConfirmSignUpResponse(result *results.ConfirmSignUpResult) *v1.ConfirmSignUpData {
	return &v1.ConfirmSignUpData{}
}

func ToLogInCommand(proto *v1.LogInRequest) *commands.LogInCommand{
	return &commands.LogInCommand{
		Email: proto.Email,
		Password: proto.Password,
	}
}
func ToLogInResponse(result *results.LogInResult) *v1.LogInData{
	return &v1.LogInData{
	}
}