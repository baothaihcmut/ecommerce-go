package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/baothaihcmut/Ecommerce-go/api-gateway/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/api-gateway/internal/middlewares"
	userProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)



type Server struct{
	cfg *config.CoreConfig
}


func NewServer(cfg *config.CoreConfig) *Server {
	return &Server{
		cfg: cfg,
	}
}

func initGrpcClient(
	sv *runtime.ServeMux,
	registerFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error,
	addr string, 
	options ...grpc.DialOption	) (*grpc.ClientConn,error) {
	conn,err := grpc.NewClient(addr,options...)
	if err != nil {
		return nil,err
	}
	if err:= registerFunc(context.Background(),sv,conn); err!= nil {
		return nil,err
	}
	return conn,nil
	
}

func (s *Server) Start() {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			return metadata.Pairs("prefix", s.cfg.Web.Prefix)
		}),
		runtime.WithForwardResponseOption(middlewares.SetHTTPStatusFromResponse),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &middlewares.CustomMarshaller{}),
		runtime.WithErrorHandler(middlewares.CustomErrorHandler),
	)
	authMux := middlewares.AuthMiddleware(s.cfg.Jwt, s.cfg.Web)(mux)
	
	clientOptions :=[]grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	userConn,err:= initGrpcClient(mux,userProto.RegisterAuthServiceHandler,s.cfg.Address.UserService,clientOptions...)
	if err != nil{
		fmt.Println(err)
		return 
	}
	defer userConn.Close()

	errCh := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errCh <- fmt.Errorf("%s", <-c)
	}()
	go func ()  {
		httpMux :=http.NewServeMux()
		httpMux.Handle(s.cfg.Web.Prefix+"/", http.StripPrefix(s.cfg.Web.Prefix, authMux))
		fmt.Println("Server is running")
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d",s.cfg.Server.Port),
			Handler: httpMux,
		}
		if err := server.ListenAndServe(); err != nil {
			errCh<-err
		}
	}()
	err=<-errCh
	fmt.Println(err)
	fmt.Println("Shut down server")	
}

