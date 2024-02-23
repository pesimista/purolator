package soap

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/pesimista/purolator-rest-api/internal/api/models"
	"github.com/pesimista/purolator-rest-api/internal/api/openapi"
)

var invalidXML string = `<s:Envelope><s:Header>`

func Test_CreateShipment(t *testing.T) {

	type args struct {
		shipment *openapi.CreateShipmentRequest
		client   HttpClient
	}

	trackinNo := "329039229987"

	shipmentCreatedXML := `<s:Envelope>
		<s:Header>
				<h:ResponseContext>
						<h:ResponseReference>ssss</h:ResponseReference>
				</h:ResponseContext>
		</s:Header>
		<s:Body>
				<CreateShipmentResponse>
						<ResponseInformation>
								<Errors/>
								<InformationalMessages i:nil="true"/>
						</ResponseInformation>
						<ShipmentPIN>
								<Value>` + trackinNo + `</Value>
						</ShipmentPIN>
						<PiecePINs>
								<PIN>
										<Value>` + trackinNo + `</Value>
								</PIN>
						</PiecePINs>
						<ReturnShipmentPINs/>
						<ExpressChequePIN>
								<Value/>
						</ExpressChequePIN>
				</CreateShipmentResponse>
		</s:Body>
	</s:Envelope>`

	errorShipmentXML := `<s:Envelope>
		<s:Header>
				<h:ResponseContext>
						<h:ResponseReference>ssss</h:ResponseReference>
				</h:ResponseContext>
		</s:Header>
		<s:Body>
				<CreateShipmentResponse>
						<ResponseInformation>
								<Errors>
										<Error>
												<Code>1100759</Code>
												<Description>Invalid Shipment Date</Description>
												<AdditionalInformation>Shipping Error</AdditionalInformation>
										</Error>
								</Errors>
								<InformationalMessages i:nil="true"/>
						</ResponseInformation>
						<ShipmentPIN i:nil="true"/>
						<PiecePINs i:nil="true"/>
						<ReturnShipmentPINs i:nil="true"/>
						<ExpressChequePIN i:nil="true"/>
				</CreateShipmentResponse>
		</s:Body>
	</s:Envelope>`

	shipment := openapi.CreateShipmentRequest{}
	err := faker.FakeData(&shipment)
	if err != nil {
		t.Fatalf("soap.CreateShipment() error = %v", err)
		return
	}

	testCases := []struct {
		name    string
		args    args
		want    *models.CreateShipmentResponse
		wantErr error
	}{
		{
			name: "When gets a valid response, return the tracking numbers",
			args: args{
				shipment: &shipment,
				client: MockHttpClient{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(shipmentCreatedXML))),
					},
					err: nil,
				},
			},
			want: &models.CreateShipmentResponse{
				PiecePINs:   []string{trackinNo},
				ShipmentPIN: trackinNo,
			},
			wantErr: nil,
		},
		{
			name: "When given an error on the response, return error",
			args: args{
				shipment: nil,
				client: MockHttpClient{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(errorShipmentXML))),
					},
					err: nil,
				},
			},
			want:    nil,
			wantErr: ErrSoapResponse,
		},
		{
			name: "When the response is an invalid XML, return error",
			args: args{
				shipment: nil,
				client: MockHttpClient{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(invalidXML))),
					},
					err: nil,
				},
			},
			want:    nil,
			wantErr: ErrInvalidXML,
		},
		{
			name: "When something goes wrong with the http request, return error",
			args: args{
				shipment: nil,
				client: MockHttpClient{
					response: &http.Response{},
					err:      nil,
				},
			},
			want:    nil,
			wantErr: ErrInvalidResponseBody,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			soapClient := NewSoapClient("something", "somekey", tt.args.client)

			got, err := soapClient.CreateShipment(tt.args.shipment)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("soap.CreateShipment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("soap.CreateShipment() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_VoidShipment(t *testing.T) {
	type args struct {
		trackingNo string
		client     HttpClient
	}

	validResponseXML := `<s:Envelope>
	<s:Header>
			<h:ResponseContext>
					<h:ResponseReference>Rating Example</h:ResponseReference>
			</h:ResponseContext>
	</s:Header>
	<s:Body>
			<VoidShipmentResponse>
					<ResponseInformation>
							<Errors/>
							<InformationalMessages i:nil="true"/>
					</ResponseInformation>
					<ShipmentVoided>true</ShipmentVoided>
			</VoidShipmentResponse>
	</s:Body>
</s:Envelope>`

	errorResponseXML := `<s:Envelope>
		<s:Header>
				<h:ResponseContext>
						<h:ResponseReference>Rating Example</h:ResponseReference>
				</h:ResponseContext>
		</s:Header>
		<s:Body>
				<VoidShipmentResponse>
						<ResponseInformation>
								<Errors>
										<Error>
												<Code>1100429</Code>
												<Description>You can only cancel shipments created today.</Description>
												<AdditionalInformation>Shipping Error</AdditionalInformation>
										</Error>
								</Errors>
								<InformationalMessages i:nil="true"/>
						</ResponseInformation>
						<ShipmentVoided>false</ShipmentVoided>
				</VoidShipmentResponse>
		</s:Body>
	</s:Envelope>`

	testCases := []struct {
		name    string
		args    args
		want    *models.VoidShipmentResponse
		wantErr error
	}{
		{
			name: "When it's a valid tracking number, return a valid response",
			args: args{
				client: MockHttpClient{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(validResponseXML))),
					},
					err: nil,
				},
				trackingNo: "329039324911",
			},
			want: &models.VoidShipmentResponse{
				ShipmentVoided: true,
			},
			wantErr: nil,
		},
		{
			name: "When the response includes an error, return an error",
			args: args{
				client: MockHttpClient{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(errorResponseXML))),
					},
					err: nil,
				},
				trackingNo: "329039324911",
			},
			want:    nil,
			wantErr: ErrSoapResponse,
		},
		{
			name: "When the response is an invalid XML, return error",
			args: args{
				trackingNo: "329039324911",
				client: MockHttpClient{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(invalidXML))),
					},
					err: nil,
				},
			},
			want:    nil,
			wantErr: ErrInvalidXML,
		},
		{
			name: "When the tracking number is missing or invalid, return error",
			args: args{
				trackingNo: "",
				client: MockHttpClient{
					response: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte(invalidXML))),
					},
					err: nil,
				},
			},
			want:    nil,
			wantErr: ErrMissingTrackingNumber,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			soapClient := NewSoapClient("client", "secret", tt.args.client)

			got, err := soapClient.VoidShipment(tt.args.trackingNo)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("soap.VoidShipment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("soap.VoidShipment() = %v, want %v", got, tt.want)
			}
		})
	}
}
