package handlers

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
)

const shipmentURL string = "https://devwebservices.purolator.com"

func (s *server) GetDocument(c *gin.Context, trackingNO string) {
	const op string = "hanlders.CreateShipment"

	xml.Marshal(op)
}
