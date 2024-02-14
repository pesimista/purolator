package models

import "encoding/xml"

type EnvelopeGetDocumentResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  struct {
		ResponseContext RequestContext
	} `xml:"Header"`
	Body GetDocumentsResponse `xml:"Body>GetDocumentsResponse"`
}

type GetDocumentsResponse struct {
	PurolatorResponseError
	Documents []DocumentInformation `xml:"Documents>Document"`
}

type DocumentInformation struct {
	TrackingNo      string `xml:"PIN>Value"`
	DocumentDetails []struct {
		DocumentType   string `xml:"DocumentType"`
		DocumentStatus string `xml:"DocumentStatus"`
		URL            string `xml:"URL"`
		Data           string `xml:"Data"`
	} `xml:"DocumentDetails>DocumentDetail"`
}
