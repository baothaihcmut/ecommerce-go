package routers

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	errHandler "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/utils"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/users/dtos"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/users/handlers"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

type UserRouter interface {
	InitRouter(*mux.Router)
}

type UserRouterImpl struct {
	handler handlers.UserHandler
	logger  *log.Logger
}

func (u *UserRouterImpl) InitRouter(r *mux.Router) {
	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.Use(middlewares.ValidateDTOMiddleware[dtos.CreateUserDto]())
	userRouter.HandleFunc("", u.handleCreateUser).Methods("POST")
}

func NewUserRouter(handler handlers.UserHandler, logger *log.Logger) UserRouter {
	return &UserRouterImpl{
		handler: handler,
		logger:  logger,
	}
}

func (u *UserRouterImpl) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(constance.PayloadContext).(*dtos.CreateUserDto)
	resp, err := u.handler.CreateUser(r.Context(), *payload)
	if err != nil {
		errHandler.HandleAppError(u.logger, w, nil, err)
	} else if !resp.Status.Success {
		errHandler.HandleAppError(u.logger, w, resp.Status, nil)
	} else {
		utils.WriteResponseSucess(w, http.StatusCreated, []string{"Create user success"}, resp)
	}
}
