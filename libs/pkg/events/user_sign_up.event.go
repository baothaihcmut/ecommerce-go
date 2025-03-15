package events

import "github.com/google/uuid"

type UserSignUpEvent struct {
	Id          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	ConfirmUrl  string    `json:"confirm-url"`
}
