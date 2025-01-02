package transports

import "github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"

type GrpcResponse struct {
	Data   interface{}
	Status *proto.Status
}
