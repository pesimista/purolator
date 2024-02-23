package soap

import (
	"testing"

	"github.com/pesimista/purolator-rest-api/internal/api/models"
)

func Test_GetDocument(t *testing.T) {
	type args struct {
		trackingNo string
		client     HttpClient
	}

	testCases := []struct {
		name    string
		args    args
		want    *models.CreateShipmentResponse
		wantErr error
	}{}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {})
	}
}
