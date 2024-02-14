package handlers

import (
	"github.com/pesimista/purolator-api/internal/api/openapi"
	"github.com/pesimista/purolator-api/internal/api/soap"
)

type server struct {
	client *soap.SoapClient
}

func NewServer(client *soap.SoapClient) openapi.ServerInterface {
	return &server{
		client: client,
	}
}
