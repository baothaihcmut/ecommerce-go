package valueobject

import (
	"regexp"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/exception"
)

var (
	phoneNumberRegex = `^\d{10}$`
)

type PhoneNumber string

func NewPhoneNumber(s string) (*PhoneNumber, error) {
	re := regexp.MustCompile(phoneNumberRegex)
	if !re.MatchString(s) {
		return nil, exception.InvalidPhonenumber
	}
	return (*PhoneNumber)(&(s)), nil
}
