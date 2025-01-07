package handlers

import (
	"net/http"

	proto "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/proto/shared"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/utils"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func HandleAppError(logger *log.Logger, w http.ResponseWriter, status *proto.Status, err error) {
	//handle internal err
	if err != nil {
		level.Error(*logger).Log("Err:", err)
		utils.WriteResponseErr(w, http.StatusInternalServerError, []string{"Internal Error"})
	}
	//handle service err
	if status != nil && !status.Success {
		utils.WriteResponseErr(w, utils.MapGrpcCodeToHttpCode(status.Code), []string{status.Message})
		return
	}

}
