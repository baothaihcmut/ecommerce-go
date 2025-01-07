package valueobject

import (
	"errors"
	"regexp"
)

type Email string

var (
	InvalidEmail = errors.New("Email not valid")
	emailRegex   = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

func NewEmail(email string) (Email, error) {
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return "", InvalidEmail
	}
	return Email(email), nil
}
