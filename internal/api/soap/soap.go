package soap

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type SoapClient struct {
	token string
}

func NewSoapClient(appKey, appSecret string) *SoapClient {
	return &SoapClient{
		token: base64.StdEncoding.EncodeToString([]byte(appKey + ":" + appSecret)),
	}
}

func (s SoapClient) HttpRequest(url, method, soapAction, body string) (string, error) {
	xmlBody := []byte(body)
	bodyReader := bytes.NewReader(xmlBody)
	req, err := http.NewRequest(method, url, bodyReader)

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return "", err
	}

	req.Header.Add("soapAction", soapAction)
	req.Header.Add("Authorization", "Basic "+s.token)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return "", err
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return "", err
	}

	return string(resBody), nil
}
