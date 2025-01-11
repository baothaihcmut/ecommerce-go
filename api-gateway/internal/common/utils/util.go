package utils

import (
	"encoding/json"
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/response"

	"google.golang.org/grpc/codes"
)

func MapGrpcCodeToHttpCode(code codes.Code) int {

	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.NotFound:
		return http.StatusNotFound
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.InvalidArgument:
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
