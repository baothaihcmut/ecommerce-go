package exception

import "errors"

var (
	ErrEmailExist       = errors.New("email exist")
	ErrPhoneNumberExist = errors.New("phone number exist")
	InvalidPhonenumber  = errors.New("Invalid phone number")
	InvalidPoint        = errors.New("Point must greater than 0")
)
