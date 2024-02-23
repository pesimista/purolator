package soap

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/go-faker/faker/v4"
)

type MockHttpClient struct {
	response *http.Response
	err      error
}

func (c MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.response, c.err
}

func Test_HttpRequest(t *testing.T) {
	type args struct {
		url        string
		method     string
		soapAction string
		body       string
	}

	testCases := []struct {
		name    string
		client  MockHttpClient
		args    args
		want    string
		wantErr error
	}{
		{
			name: "When called, return response value",
			client: MockHttpClient{
				response: &http.Response{Body: io.NopCloser(bytes.NewReader([]byte("<Envelope></Envelope>")))},
				err:      nil,
			},
			args: args{
				url:        shippingServiceURL,
				method:     http.MethodPost,
				soapAction: createShipmentAction,
				body:       "<soap:Envelope></soap:Envelope>",
			},
			want:    "<Envelope></Envelope>",
			wantErr: nil,
		},
		{
			name: "When called with invalid url, return error",
			client: MockHttpClient{
				response: &http.Response{},
				err:      nil,
			},
			args: args{
				url:        "\n/api",
				method:     http.MethodPost,
				soapAction: createShipmentAction,
				body:       "<soap:Envelope></soap:Envelope>",
			},
			want:    "",
			wantErr: ErrInvalidRequestURL,
		},
		{
			name: "When called and gets an invalid body, return error",
			client: MockHttpClient{
				response: &http.Response{},
				err:      nil,
			},
			args: args{
				url:        shippingServiceURL,
				method:     http.MethodPost,
				soapAction: createShipmentAction,
				body:       "<soap:Envelope></soap:Envelope>",
			},
			want:    "",
			wantErr: ErrInvalidResponseBody,
		},
		{
			name: "When gets an error while executing the request, return error",
			client: MockHttpClient{
				response: &http.Response{
					Body: io.NopCloser(bytes.NewReader([]byte("response"))),
				},
				err: fmt.Errorf("failed"),
			},
			args: args{
				url:        shippingServiceURL,
				method:     http.MethodPost,
				soapAction: createShipmentAction,
				body:       "<soap:Envelope></soap:Envelope>",
			},
			want:    "",
			wantErr: ErrFailedRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			soapClient := NewSoapClient("key", "pass", tt.client)

			got, err := soapClient.HttpRequest(tt.args.url, tt.args.method, tt.args.soapAction, tt.args.body)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("soap.HttpRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Fatalf("soap.HttpRequest() %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func Test_NewEnvelopeXML(t *testing.T) {
	type MockXMLBody struct {
		Pin string `xml:"PIN>Value"`
	}

	pin := faker.UUIDDigit()

	testCases := []struct {
		name    string
		args    MockXMLBody
		want    string
		wantErr error
	}{
		{
			name: "When called with valid body, return xml as string",
			args: MockXMLBody{
				Pin: pin,
			},
			want: `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:q2="http://purolator.com/pws/datatypes/v2">` +
				`<soap:Header>` +
				`<q2:RequestContext>` +
				`<q2:Version>2.0</q2:Version>` +
				`<q2:Language>en</q2:Language>` +
				`<q2:GroupID>234521</q2:GroupID>` +
				`<q2:RequestReference>uuid</q2:RequestReference>` +
				`</q2:RequestContext>` +
				`</soap:Header>` +
				`<soap:Body>` +
				`<q2:MockXMLBody>` +
				`<q2:PIN>` +
				`<q2:Value>` + pin + `</q2:Value>` +
				`</q2:PIN>` +
				`</q2:MockXMLBody>` +
				`</soap:Body>` +
				`</soap:Envelope>`,
			wantErr: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEnvelopeXML(tt.args)
			uuidRegex := regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
			regex := regexp.MustCompile(`\n\s*`)

			got = uuidRegex.ReplaceAllString(got, "uuid") // just remove the uuid
			got = regex.ReplaceAllString(got, "")         // remove spaces

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("soap.NewEnvelopeXML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != got {
				t.Fatalf("soap.NewEnvelopeXML() = %v, want %v", got, tt.want)
			}
		})
	}
}
