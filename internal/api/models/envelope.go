package models

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	bodyPrefix   = "q2"
	headerPrefix = "soap"
)

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

func NewEnvelopeXML(body any) (string, error) {
	envelope := &Envelope{
		Soap:   "http://schemas.xmlsoap.org/soap/envelope/",
		Q2:     "http://purolator.com/pws/datatypes/v2",
		Header: "REPLACED_BY_HEADER",
		Body:   "REPLACED_BY_BODY",
	}

	envelopeBuffer, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return "", err
	}
	envelopeXML := addPrefix(string(envelopeBuffer), headerPrefix)

	var (
		Version          = "2.0"
		Language         = "en"
		GroupID          = "234521"
		RequestReference = uuid.New().String()
	)

	header := RequestContext{
		Version:          &Version,
		Language:         &Language,
		GroupID:          &GroupID,
		RequestReference: &RequestReference,
	}

	headerBuffer, err := xml.MarshalIndent(header, "", "  ")
	if err != nil {
		return "", err
	}
	headerXML := addPrefix(string(headerBuffer), bodyPrefix)

	bodyBuffer, err := xml.MarshalIndent(body, "", "  ")
	if err != nil {
		return "", err
	}
	bodyXML := addPrefix(string(bodyBuffer), bodyPrefix)

	envelopeXML = strings.ReplaceAll(envelopeXML, "REPLACED_BY_HEADER", fmt.Sprintf("\n%s\n", headerXML))
	envelopeXML = strings.ReplaceAll(envelopeXML, "REPLACED_BY_BODY", fmt.Sprintf("\n%s\n", bodyXML))

	fmt.Println(envelopeXML)
	return envelopeXML, nil
}

func addPrefix(data, prefix string) string {
	regex := regexp.MustCompile(`(<\/?)`)

	return regex.ReplaceAllString(data, fmt.Sprintf("${1}%s:", prefix))
}
