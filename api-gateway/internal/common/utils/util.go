package utils

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/proto"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/models"
	"google.golang.org/grpc/codes"
)

func MapGrpcCodeToHttpCode(code string) int {

	switch code {
	case codes.OK.String():
		return http.StatusOK
	case codes.AlreadyExists.String():
		return http.StatusConflict
	case codes.NotFound.String():
		return http.StatusNotFound
	case codes.PermissionDenied.String():
		return http.StatusForbidden
	case codes.Unauthenticated.String():
		return http.StatusUnauthorized
	case codes.InvalidArgument.String():
		return http.StatusBadRequest
	default:
		return 500
	}
}
func MapRole(src models.Role) proto.Role {
	switch src {
	case models.RoleCustomer:
		return proto.Role_CUSTOMER
	case models.RoleShopOwner:
		return proto.Role_SHOP_OWNER
	case models.RoleAdmin:
		return proto.Role_ADMIN
	}
	return proto.Role_CUSTOMER
}
func MapRoleProto(src proto.Role) models.Role {
	switch src {
	case proto.Role_CUSTOMER:
		return models.RoleCustomer
	case proto.Role_SHOP_OWNER:
		return models.RoleShopOwner
	case proto.Role_ADMIN:
		return models.RoleAdmin
	default:
		return models.RoleCustomer
	}
}
