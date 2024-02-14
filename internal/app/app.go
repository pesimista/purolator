package app

import (
	"github.com/pesimista/purolator-api/internal/api"
)

func Run() {
	server := api.NewServer()
	// server.CreateServer()

	server.SetRoutes()
	server.Run()
}
