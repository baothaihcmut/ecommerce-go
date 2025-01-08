package valueobject

import (
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
)

var (
	ErrTokenTypeMisMatch  = errors.New("Token type mismatch")
	ErrTokenValueMisMatch = errors.New("Token value mismatch")
)

type Token struct {
	Value     string
	TokenType enums.TokenType
}

func (r Token) IsEqual(token Token) bool {

	return r.Value == token.Value && r.TokenType == token.TokenType
}
