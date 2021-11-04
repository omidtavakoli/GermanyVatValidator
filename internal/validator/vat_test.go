package validator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestHandler_Vat(t *testing.T) {
	tests := []struct {
		description        string
		expectedStatusCode int
		input              string
		expectedData       bool
		initService        func() Service
	}{
		{
			description:        "200 ok",
			expectedStatusCode: http.StatusOK,
			input:              "PL9372717673",
			expectedData:       true,
			initService: func() Service {
				mockRep := new(MockProxyService)
				mockRep.On("VatValidator", mock.Anything).Return("ok", nil)
				return mockRep
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			service := tt.initService()
			ret, err := service.VatValidator(tt.input)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, ret, tt.expectedData)
		})
	}
}
