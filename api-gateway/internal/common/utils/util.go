package utils

import (
	"net/http"

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
