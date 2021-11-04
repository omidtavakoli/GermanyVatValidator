package EU_VIES

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCheckVAT(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 10000 * time.Millisecond,
	}
	appClient := client{
		HttpClient: httpClient,
		Config:     &Config{Url: "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"},
	}

	tests := []struct {
		description  string
		expectedData bool
		input        string
	}{
		{
			description:  "Home24, valid",
			expectedData: true,
			input:        "DE266182271",
		},
		{
			description:  "Bosch, valid",
			expectedData: true,
			input:        "DE811128135",
		},
		{
			description:  "invalid",
			expectedData: false,
			input:        "DE948734627",
		},
		{
			description:  "valid",
			expectedData: true,
			input:        "IE6388047V",
		},
		{
			description:  "valid",
			expectedData: true,
			input:        "PL9372717673",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			vatData, err := appClient.CheckVAT(tt.input)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, vatData.Valid, tt.expectedData)
			time.Sleep(2 * time.Second) // prevent from banning by external endpoint
		})
	}

}
