package valueobject

import "golang.org/x/crypto/bcrypt"

type Password string

func NewPassword(password string) (Password, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return Password(hashedPassword), nil
}

func (p Password) Compare(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(string(p)), []byte(password))
	return err == nil
}
