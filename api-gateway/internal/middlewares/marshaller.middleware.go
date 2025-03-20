package middlewares

import (
	"encoding/json"
	"reflect"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type CustomMarshaller struct {
    runtime.JSONPb
}


// Marshal transforms the gRPC response before sending it to the client
func (c *CustomMarshaller) Marshal(v any) ([]byte, error) {
	var message string
	var data any

    vReflect := reflect.Indirect(reflect.ValueOf(v))
	statusField := vReflect.FieldByName("Status")
	if statusField.IsValid() && statusField.Kind() == reflect.Ptr {
		statusField = statusField.Elem()
		messageField := statusField.FieldByName("Message")
		if messageField.IsValid() && messageField.Kind() == reflect.String {
			message = messageField.String()
		}
	}
    dataField := vReflect.FieldByName("Data")
	if dataField.IsValid() {
		if dataField.Kind() == reflect.Ptr && !dataField.IsNil() {
			dataField = dataField.Elem()
		}

		data = dataField.Interface()
	}
    // Build the transformed response
    modifiedResponse := map[string]any{
        "success": true,
        "message": message,
        "data":    data,
    }

    return json.Marshal(modifiedResponse)
}
