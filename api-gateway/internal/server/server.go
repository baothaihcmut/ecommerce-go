package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/discovery"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
)

type Server struct {
	ConsulClient *api.Client
	Cfg          *config.Config
	Logger       *log.Logger
}

func NewServer(consul *api.Client, cfg *config.Config, logger *log.Logger) *Server {
	return &Server{
		ConsulClient: consul,
		Cfg:          cfg,
		Logger:       logger,
	}
}

func (s *Server) Run() {
	//init Service
	discoveryService := discovery.NewDiscoveryService(s.ConsulClient)

	//add middleware
	r := mux.NewRouter()
	r.Use(handlers.CORS())

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Cfg.ServerConfig.Port), r); err != nil {
			level.Error(*s.Logger).Log("err", "Error start http server")
			panic(err)
		}
		level.Info(*s.Logger).Log("msg", "Server started successfully 🚀")
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdown()
	<-ctx.Done()
	level.Info(*s.Logger).Log("msg", "Server gracefully stopped.")

}
