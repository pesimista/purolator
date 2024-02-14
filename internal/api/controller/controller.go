package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/pesimista/purolator-api/internal/api/handlers"
	"github.com/pesimista/purolator-api/internal/api/openapi"
	"github.com/pesimista/purolator-api/internal/api/soap"
)

func NewRouter(handler *gin.Engine) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	// handler.Use(middleware.())

	handler.StaticFile("/swagger", "./spec/openapi.yaml")
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger", Path: "/swagger/ui"}
	sh := middleware.SwaggerUI(opts, nil)
	handler.GET("/swagger/ui", func(ctx *gin.Context) {
		sh.ServeHTTP(ctx.Writer, ctx.Request)
	})

	opt := openapi.GinServerOptions{
		BaseURL:     "/api/v1",
		Middlewares: make([]openapi.MiddlewareFunc, 0),
	}

	client := soap.NewSoapClient("", "")
	RegisterHandlers(handler, handlers.NewServer(client), opt)
}

func RegisterHandlers(router *gin.Engine, si openapi.ServerInterface, options openapi.GinServerOptions) *gin.Engine {
	wrapper := openapi.ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	router.GET(options.BaseURL+"/shipments", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"msg": "hello"}) })
	router.POST(options.BaseURL+"/shipments", wrapper.CreateShipment)
	router.DELETE(options.BaseURL+"/shipments/:trackingNo", wrapper.VoidShipment)

	return router
}
