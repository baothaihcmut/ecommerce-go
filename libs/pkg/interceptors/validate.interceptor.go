package interceptors

import (
	"context"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	reqVal := reflect.ValueOf(req)
	validateMethod := reqVal.MethodByName("Validate")
	if validateMethod.IsValid() {
		res := validateMethod.Call(nil)
		if len(res) > 0 && !res[0].IsNil() {
			if err, ok := res[0].Interface().(error); ok {
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
		}
	}
	return handler(ctx, req)
}
