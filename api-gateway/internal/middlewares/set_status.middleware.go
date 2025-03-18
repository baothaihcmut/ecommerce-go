package middlewares

import (
	"context"
	"net/http"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// SetHTTPStatusFromResponse extracts the status code from the response and sets the HTTP status.
func SetHTTPStatusFromResponse(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
    v := reflect.Indirect(reflect.ValueOf(resp))
    statusField := v.FieldByName("Status")
    if statusField.IsValid() && statusField.Kind() == reflect.Ptr {
        codeField := statusField.Elem().FieldByName("Code")
        if codeField.IsValid() && codeField.CanInt() {
            httpStatus := int(codeField.Int())
            w.WriteHeader(httpStatus) // 
        }
    }

    return nil
}