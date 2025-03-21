package interceptors

import (
	"context"
	"slices"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/constant"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/models"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ExtractAuthInterceptor(publicServices ...string) func(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if slices.Contains(publicServices, info.FullMethod) {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Missing auth info")
		}
		userId, exist := md["user_id"]
		if !exist || len(userId) == 0 {
			return nil, status.Error(codes.Unauthenticated, "Missing auth info")
		}
		isShopOwnerActive, exist := md["is_shop_owner_active"]
		if !exist || len(isShopOwnerActive) == 0 {
			return nil, status.Error(codes.Unauthenticated, "Missing auth info")
		}
		uUserId, _ := uuid.FromBytes([]byte(userId[0]))
		bisShopOwnerActive := false
		if isShopOwnerActive[0] == "true" {
			bisShopOwnerActive = true
		}
		userCtx := models.UserContext{
			UserId:            uUserId,
			IsShopOwnerActive: bisShopOwnerActive,
		}
		ctx = context.WithValue(ctx, constant.UserContextKey, &userCtx)
		return handler(ctx, req)
	}
}
