package exception

import "errors"

var (
	ErrEmailExist            = errors.New("email exist")
	ErrPhonenumberExist      = errors.New("phone number exist")
	ErrUserPendingForConfirm = errors.New("user is pending for confirm")
)
