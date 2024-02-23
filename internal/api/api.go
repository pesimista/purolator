//go:generate oapi-codegen -generate gin -o openapi/openapi_gin.gen.go -package openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -generate types -o openapi/openapi_types.gen.go -package openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -generate client -o openapi/openapi_client.gen.go -package openapi ../../spec/openapi.yaml
//go:generate oapi-codegen -generate spec -o openapi/openapi_spec.gen.go -package openapi ../../spec/openapi.yaml

package api

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pesimista/purolator-rest-api/internal/api/controller"
)

type Server struct {
	engine *gin.Engine
	server *http.Server
}

func NewServer() *Server {
	handler := gin.New()
	return &Server{
		engine: handler,
		server: &http.Server{
			Addr:              "localhost:8080",
			Handler:           handler,
			ReadHeaderTimeout: time.Second * 30,
		},
	}
}

// func (s *Server) CreateServer() {
// if err:=sentry.Config();err!=nil {
// }

// 	s.engine.NoRoute(func(ctx *gin.Context) {
// 		fmt.Println("Route not found")
// 	})
// }

func (s *Server) SetRoutes() {
	controller.NewRouter(s.engine)
}

func (s *Server) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s", err)
		}
	}()

	fmt.Println("App Started")

	<-ctx.Done()

	stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Printf("listen: %s", err)
	}

	fmt.Println("App Exiting")

}
