package soap

import (
	"encoding/xml"
)

func (s *SoapClient) GetDocument(trackingNO string) {
	const op string = "hanlders.CreateShipment"

	xml.Marshal(op)
}
