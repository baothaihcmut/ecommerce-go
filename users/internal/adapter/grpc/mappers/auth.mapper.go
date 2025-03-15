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
