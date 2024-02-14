package models

import "encoding/xml"

type EnvelopeCreateShipmentResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  struct {
		ResponseContext RequestContext
	} `xml:"Header"`
	Body CreateShipmentResponse `xml:"Body>CreateShipmentResponse"`
}

type CreateShipmentResponse struct {
	PurolatorResponseError

	ShipmentPIN string `xml:"ShipmentPIN>Value" json:"masterTrackingId,omitempty"`

	PiecePINs          []string `xml:"PiecePINs>PIN>Value" json:"trackingNumbers,omitempty"`
	ReturnShipmentPINs []string `xml:"ReturnShipmentPINs>PIN>Value" json:"returnTrackingNumber,omitempty"`
	ExpressChequePIN   []string `xml:"ExpressChequePIN>PIN>Value" json:"expressChequePIN,omitempty"`
}

type VoidShipmentRequest struct {
	Pin string `xml:"PIN>Value" json:"trackingNumber"`
}

type EnvelopeVoidShipmentResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  struct {
		ResponseContext RequestContext
	} `xml:"Header"`
	Body VoidShipmentResponse `xml:"Body>VoidShipmentResponse"`
}

type VoidShipmentResponse struct {
	PurolatorResponseError
	ShipmentVoided bool `xml:"ShipmentVoided" json:"shipmentVoided"`
}
