package exception

import (
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/exception"
	"google.golang.org/grpc/codes"
)

var errorMap map[error]codes.Code = map[error]codes.Code{
	exception.ErrCategoryNotExist:       codes.NotFound,
	exception.ErrShopNotExist:           codes.NotFound,
	exception.ErrUserNotShopOwnerActive: codes.PermissionDenied,
	exception.ErrUserIsNotShopOwner:     codes.PermissionDenied,
}

func MapException(err error) (codes.Code, string) {
	if code, exist := errorMap[err]; exist {
		return code, err.Error()
	}
	return codes.Internal, "Internal server error"
}
