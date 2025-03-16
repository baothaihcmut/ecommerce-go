package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/events"

	"github.com/baothaihcmut/Ecommerce-go/mail/internal/mailer"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/models"
)
type AuthMailService interface {
	SendMailConfirmSignUp(ctx context.Context, e *events.UserSignUpEvent) error
	
}
type AuthMailServiceImpl struct {
	mailer mailer.Mailer
}
func (a *AuthMailServiceImpl) SendMailConfirmSignUp(ctx context.Context, e *events.UserSignUpEvent) error {
	err := a.mailer.SendMail(
		ctx,
		mailer.SendMailArg{
			Subject:  "Sign Up verfication",
			To:       e.Email,
			Template: "auth/confirm_sign_up.html",
			Data: models.ConfirmSignUpModel{
				Email:      e.Email,
				LastName:   e.LastName,
				FirstName:  e.FirstName,
				ConfirmUrl: e.ConfirmUrl,
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func NewAuthMailService(mailer mailer.Mailer) AuthMailService {
	return &AuthMailServiceImpl{
		mailer: mailer,
	}
}
