package handlers

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/utils"
	"github.com/go-kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcErrorHandler(logger *log.Logger, w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		utils.WriteResponseErr(w, http.StatusInternalServerError, []string{"Internal err"})
	}
	code := st.Code()
	msg := st.Message()
	if code != codes.OK {
		utils.WriteResponseErr(w, utils.MapGrpcCodeToHttpCode(code), []string{msg})
	}
}
