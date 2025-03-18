package exception

import (
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/exception"
	"google.golang.org/grpc/codes"
)



var errorMap map[error]codes.Code = map[error]codes.Code{
	exception.ErrEmailExist:            codes.AlreadyExists,
	exception.ErrPhonenumberExist:      codes.AlreadyExists,
	exception.ErrUserPendingForConfirm: codes.AlreadyExists,
	exception.ErrWrongEmailOrPassword: codes.Unauthenticated,
}

func MapException(err error) ( codes.Code,string) {
	if code, exist := errorMap[err]; exist {
		return code,err.Error()
	}
	return  codes.Internal,"Internal server error"
}
