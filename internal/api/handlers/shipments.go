package handlers

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	cErrors "github.com/pesimista/purolator-rest-api/internal/api/errors"
	"github.com/pesimista/purolator-rest-api/internal/api/openapi"
)

const (
	billingAccount string = "9999999999"
)

func (s *server) CreateShipment(c *gin.Context) {
	const op string = "hanlders.CreateShipment"

	var shipment *openapi.CreateShipmentRequest
	if err := c.ShouldBindJSON(&shipment); err != nil {
		cErrors.JSON(c, op, "could not bind request body", err, http.StatusInternalServerError)
		return
	}

	account := billingAccount

	shipment.Shipment.PaymentInformation.RegisteredAccountNumber = &account
	shipment.Shipment.PaymentInformation.BillingAccountNumber = &account

	data, err := s.client.CreateShipment(shipment)
	if err != nil {
		cErrors.JSON(c, op, "", err, http.StatusInternalServerError)
		return
	}

	if data.Error != nil {
		cErrors.JSON(c, op, data.Error.Description, err, http.StatusBadRequest)
		return
	}

	bytesaar, _ := json.MarshalIndent(data, "", "  ")

	fmt.Println(string(bytesaar))

	c.JSON(
		http.StatusCreated,
		openapi.CreateShipmentRes{
			MasterTrackingNo: data.ShipmentPIN,
			TrackingNOs:      data.PiecePINs,
		},
	)
}

func (s *server) VoidShipment(c *gin.Context, trackingNo string) {
	const op string = "hanlders.VoidShipment"

	if len(trackingNo) == 0 {
		cErrors.JSON(c, op, "missing tracking number", nil, http.StatusBadRequest)
		return
	}

	_, err := s.client.VoidShipment(trackingNo)
	if err != nil {
		cErrors.JSON(c, op, "", err, http.StatusBadRequest)
		return
	}

	c.Status(http.StatusNoContent)
}
