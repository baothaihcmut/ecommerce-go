package grpc

import (
	"time"

	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type ConnectionArgs struct {
	MaxConnection         int
	MinConnection         int
	ConnectionIdleTimeOut time.Duration
}

func GetConnectionPool(addr string, args ConnectionArgs, opts ...grpc.DialOption) (*grpcpool.Pool, error) {
	return grpcpool.New(func() (*grpc.ClientConn, error) {
		conn, err := grpc.NewClient(addr, opts...)
		if err != nil {
			return nil, err
		}
		return conn, nil
	},
		args.MaxConnection,
		args.MinConnection,
		args.ConnectionIdleTimeOut,
	)
}
