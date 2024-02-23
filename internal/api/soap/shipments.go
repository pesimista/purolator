package soap

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/pesimista/purolator-rest-api/internal/api/models"
	"github.com/pesimista/purolator-rest-api/internal/api/openapi"
)

const (
	shippingServiceURL   = "https://devwebservices.purolator.com/EWS/v2/Shipping/ShippingService.asmx"
	createShipmentAction = "http://purolator.com/pws/service/v2/CreateShipment"
	voidShipmentAction   = "http://purolator.com/pws/service/v2/VoidShipment"
)

func (s *SoapClient) CreateShipment(shipment *openapi.CreateShipmentRequest) (*models.CreateShipmentResponse, error) {
	const op string = "soap.CreateShipment"

	envelopeXML, err := NewEnvelopeXML(shipment)
	if err != nil {
		return nil, fmt.Errorf("%s: could create an envelope for the request: %s", op, err)
	}

	responseString, err := s.HttpRequest(
		shippingServiceURL,
		http.MethodPost,
		createShipmentAction,
		envelopeXML,
	)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", op, err)
	}

	var response *models.EnvelopeCreateShipmentResponse
	err = xml.Unmarshal([]byte(responseString), &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w %w", op, ErrInvalidXML, err)
	}

	if response.Body.Error != nil {
		return nil, fmt.Errorf("%s: %w %v", op, ErrSoapResponse, response.Body.Error.Description)
	}

	return &response.Body, nil
}

func (s *SoapClient) VoidShipment(trackingNo string) (*models.VoidShipmentResponse, error) {
	const op string = "soap.VoidShipment"

	if len(trackingNo) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrMissingTrackingNumber)
	}

	voidRequest := models.VoidShipmentRequest{
		Pin: trackingNo,
	}

	envelopeXML, err := NewEnvelopeXML(voidRequest)
	if err != nil {
		return nil, fmt.Errorf("%s: could create an envelope for the request: %s", op, err)
	}

	responseString, err := s.HttpRequest(
		shippingServiceURL,
		http.MethodPost,
		voidShipmentAction,
		envelopeXML,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: could create an envelope for the request: %s", op, err)
	}

	var response *models.EnvelopeVoidShipmentResponse
	err = xml.Unmarshal([]byte(responseString), &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w %s", op, ErrInvalidXML, err)
	}

	if response.Body.Error != nil {
		return nil, fmt.Errorf("%s: %w %v", op, ErrSoapResponse, response.Body.Error.Description)
	}

	return &response.Body, nil

}
