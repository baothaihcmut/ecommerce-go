package utils

import (
	"encoding/json"
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/users/dtos"
	proto "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/users/proto/users"
	"google.golang.org/grpc/codes"
)

func MapRole(role dtos.Role) proto.Role {
	switch role {
	case dtos.ADMIN:
		return proto.Role_ADMIN
	case dtos.USER:
		return proto.Role_CUSTOMER
	case dtos.SHOP_OWNER:
		return proto.Role_SHOP_OWNER
	default:
		return proto.Role_CUSTOMER
	}
}

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

func WriteResponseErr(w http.ResponseWriter, code int, message []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response.InitResponse[any](false, message, nil))

}

func WriteResponseSucess[T any](w http.ResponseWriter, code int, message []string, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response.InitResponse[T](false, message, data))
}
