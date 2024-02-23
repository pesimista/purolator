package soap

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/pesimista/purolator-rest-api/internal/api/models"
)

const (
	bodyPrefix   = "q2"
	headerPrefix = "soap"
)

var (
	ErrMissingTrackingNumber = errors.New("missing tracking number")
	ErrInvalidRequestURL     = errors.New("invalid request url")
	ErrInvalidRequestBody    = errors.New("invalid request body")
	ErrFailedRequest         = errors.New("error making http request")
	ErrInvalidResponseBody   = errors.New("could not read response body")
	ErrInvalidXML            = errors.New("could not decode xml body")
	ErrSoapResponse          = errors.New("error on soap response")
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type SoapClient struct {
	token      string
	httpClient HttpClient
}

func NewSoapClient(appKey, appSecret string, httpClient HttpClient) *SoapClient {
	return &SoapClient{
		token:      base64.StdEncoding.EncodeToString([]byte(appKey + ":" + appSecret)),
		httpClient: httpClient,
	}
}

func (s SoapClient) HttpRequest(url, method, soapAction, body string) (string, error) {
	op := "soap.HttpRequest"

	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return "", fmt.Errorf("%v: %w %w", op, ErrInvalidRequestURL, err)
	}

	req.Header.Add("soapAction", soapAction)
	req.Header.Add("Authorization", "Basic "+s.token)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")

	response, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%v: %w %w", op, ErrFailedRequest, err)
	}

	if response == nil || response.Body == nil {
		return "", fmt.Errorf("%v: %w", op, ErrInvalidResponseBody)
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("%v: %w %w", op, ErrInvalidResponseBody, err)
	}

	return string(resBody), nil
}

func NewEnvelopeXML(body any) (string, error) {
	op := "soap.NewEnvelopeXML"

	envelope := models.NewEnvelope("REPLACED_BY_HEADER", "REPLACED_BY_BODY")

	envelopeBuffer, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%v: %w %w", op, ErrInvalidRequestBody, err)
	}

	envelopeXML := addPrefix(string(envelopeBuffer), headerPrefix)

	var (
		Version          = "2.0"
		Language         = "en"
		GroupID          = "234521"
		RequestReference = uuid.New().String()
	)

	header := models.RequestContext{
		Version:          &Version,
		Language:         &Language,
		GroupID:          &GroupID,
		RequestReference: &RequestReference,
	}

	headerBuffer, err := xml.MarshalIndent(header, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%v: %w %w", op, ErrInvalidRequestBody, err)
	}
	headerXML := addPrefix(string(headerBuffer), bodyPrefix)

	bodyBuffer, err := xml.MarshalIndent(body, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%v: %w %w", op, ErrInvalidRequestBody, err)
	}
	bodyXML := addPrefix(string(bodyBuffer), bodyPrefix)

	envelopeXML = strings.ReplaceAll(envelopeXML, "REPLACED_BY_HEADER", fmt.Sprintf("\n%s\n", headerXML))
	envelopeXML = strings.ReplaceAll(envelopeXML, "REPLACED_BY_BODY", fmt.Sprintf("\n%s\n", bodyXML))

	return envelopeXML, nil
}

func addPrefix(data, prefix string) string {
	regex := regexp.MustCompile(`(<\/?)`)

	return regex.ReplaceAllString(data, fmt.Sprintf("${1}%s:", prefix))
}
