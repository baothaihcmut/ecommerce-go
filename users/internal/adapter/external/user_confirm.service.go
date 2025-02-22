package external

import (
	"context"
	"fmt"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/cache"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/queue"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/events"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type UserConfirmService struct {
	redisService cache.CacheService
	queueService queue.QueueService
	tracer       trace.Tracer
}

type UserMessage struct {
	*events.UserSignUpEvent
	Action string `json:"action"`
}

func (u *UserConfirmService) StoreUserInfo(ctx context.Context, user *user.User) (code string, err error) {
	ctx, span := tracing.StartSpan(ctx, u.tracer, "User.ConfirmService: external service", nil)
	defer tracing.EndSpan(span, err, nil)
	//generate code
	codeVerification := uuid.New().String()
	//store user info to redis
	if err = u.redisService.SetValue(ctx, fmt.Sprintf("user_verification:%s", code), *user, 30*time.Minute); err != nil {
		return "", err
	}
	//store email to pending
	if err = u.redisService.SetString(ctx, fmt.Sprintf("email_pending_verification:%s", string(user.Email)), string(user.Email), 30*time.Minute); err != nil {
		return "", err
	}
	return codeVerification, nil
}

func (u *UserConfirmService) PublishSignUpEvent(ctx context.Context, e *events.UserSignUpEvent) (err error) {
	ctx, span := tracing.StartSpan(ctx, u.tracer, "User.ConfirmService: external service", nil)
	defer tracing.EndSpan(span, err, nil)
	_, _, err = u.queueService.PublishMessage("users", UserMessage{
		UserSignUpEvent: e,
		Action:          "sign_up",
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
