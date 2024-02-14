package soap

import (
	"encoding/xml"
)

const shipmentURL string = "https://devwebservices.purolator.com"

func (s *SoapClient) GetDocument(trackingNO string) {
	const op string = "hanlders.CreateShipment"

	xml.Marshal(op)
}
