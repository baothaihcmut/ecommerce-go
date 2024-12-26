package valueobject

import (
	"errors"
	"regexp"
)

var (
	InvalidPhonenumber = errors.New("Invalid phone number")
	phoneNumberRegex   = `^\d{10}$`
)

type PhoneNumber string

func NewPhoneNumber(s string) (*PhoneNumber, error) {
	re := regexp.MustCompile(phoneNumberRegex)
	if !re.MatchString(s) {
		return nil, InvalidPhonenumber
	}
	return (*PhoneNumber)(&(s)), nil
}
