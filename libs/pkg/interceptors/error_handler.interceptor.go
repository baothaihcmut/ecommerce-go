package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)




func ErrorHandler(mapErrorFunc func(error)(codes.Code,string)) func (
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error)  {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler) (any, error) {
		res,err := handler(ctx,req)
		if err!= nil{
			code,msg:= mapErrorFunc(err)
			return nil, status.Error(code,msg)
		}
		return res,nil
	}
}

