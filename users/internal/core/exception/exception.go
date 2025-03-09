package exception

import "errors"

var (
	ErrEmailExist       = errors.New("email exist")
	ErrPhonenumberExist = errors.New("phone number exist")
)
