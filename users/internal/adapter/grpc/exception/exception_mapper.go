package exception

import (
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/exception"
	"google.golang.org/grpc/codes"
)

type errorDetail struct {
	Msg  string
	Code codes.Code
}

var errorMap map[error]codes.Code = map[error]codes.Code{
	exception.ErrEmailExist:            codes.AlreadyExists,
	exception.ErrPhonenumberExist:      codes.AlreadyExists,
	exception.ErrUserPendingForConfirm: codes.AlreadyExists,
}

func MapException(err error) (string, codes.Code) {
	if code, exist := errorMap[err]; exist {
		return err.Error(), code
	}
	return "Internal server error", codes.Internal
}
