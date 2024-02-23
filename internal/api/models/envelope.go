package models

type PurolatorResponseError struct {
	Error *struct {
		Code                  string `xml:"Code" json:"code,omitempty"`
		Description           string `xml:"Description" json:"description,omitempty"`
		AdditionalInformation string `xml:"AdditionalInformation" json:"additionalInformation,omitempty"`
	} `xml:"ResponseInformation>Errors>Error,omitempty" json:"error,omitempty"`
}

type RequestContext struct {
	Version           *string `xml:"Version,omitempty"`
	Language          *string `xml:"Language,omitempty"`
	GroupID           *string `xml:"GroupID,omitempty"`
	RequestReference  *string `xml:"RequestReference,omitempty"`
	ResponseReference *string `xml:"ResponseReference,omitempty"`
}

type Envelope struct {
	Soap   string `xml:"xmlns:soap,attr"`
	Q2     string `xml:"xmlns:q2,attr"`
	Header any
	Body   any
}

func NewEnvelope(header, body any) *Envelope {
	return &Envelope{
		Soap:   "http://schemas.xmlsoap.org/soap/envelope/",
		Q2:     "http://purolator.com/pws/datatypes/v2",
		Header: header,
		Body:   body,
	}
}
