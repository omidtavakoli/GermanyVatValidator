package rest

import (
	"VatIdValidator/internal/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var dummyTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

func TestHandler_VAT(t *testing.T) {
	tests := []struct {
		description        string
		expectedStatusCode int
		expectedData       string
		input              string
		initService        func() validator.Service
	}{
		{
			description:        "200 ok",
			expectedStatusCode: http.StatusOK,
			expectedData:       "ok",
			input:              "DE123456789",
			initService: func() validator.Service {
				mockRep := new(MockProxyService)
				mockRep.On("VatValidator", mock.Anything).Return("ok", nil)
				return mockRep
			},
		},
		{
			description:        "417 expectation failed",
			expectedStatusCode: http.StatusExpectationFailed,
			expectedData:       "expectation failed",
			input:              "De123456789",
			initService: func() validator.Service {
				mockRep := new(MockProxyService)
				mockRep.On("VatValidator", mock.Anything).Return("ok", nil)
				return mockRep
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			service := tt.initService()
			handler := CreateHandler(service)
			gin.SetMode(gin.TestMode)
			gin.DefaultWriter = ioutil.Discard
			router := gin.Default()
			router.GET("/validator/vat/:id", handler.VatValidator)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/validator/vat/%s", tt.input), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}
