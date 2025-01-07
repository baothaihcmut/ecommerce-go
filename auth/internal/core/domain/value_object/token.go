package valueobject

import (
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain/enums"
)

var (
	ErrTokenTypeMisMatch  = errors.New("Token type mismatch")
	ErrTokenValueMisMatch = errors.New("Token value mismatch")
)

type Token struct {
	Value     string
	TokenType enums.TokenType
}

func (r Token) IsEqual(token Token) error {
	if r.TokenType != token.TokenType {
		return ErrTokenTypeMisMatch
	}
	if r.Value != token.Value {
		return ErrTokenValueMisMatch
	}
	return nil
}
