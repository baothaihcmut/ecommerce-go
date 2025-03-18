// CustomErrorHandler transforms gRPC errors into JSON responses
package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type ErrorResponse struct {
	
    Success bool   `json:"success"`
    Message string `json:"message"`
}


func CustomErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
    st, _ := status.FromError(err)
    httpCode := runtime.HTTPStatusFromCode(st.Code())

    // Construct error response
    errorResponse := ErrorResponse{
        Success: false,
        Message: st.Message(),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(httpCode)

    json.NewEncoder(w).Encode(errorResponse)
}
