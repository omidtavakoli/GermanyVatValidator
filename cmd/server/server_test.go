package main

import (
	"VatIdValidator/internal/logger"
	"VatIdValidator/internal/validator"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	cfg = MainConfig{
		Logger: logger.Config{},
		Proxy:  validator.Config{},
		Server: ServerConfig{},
	}
	ctx = context.Background()
)

func TestServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = ioutil.Discard
	server := NewServer(&cfg, nil)
	err := server.Initialize(ctx)
	if err != nil {
		t.Error(err)
	}
	router := SetupRouter(server.RESTHandler, server.Config)
	ts := httptest.NewServer(router)
	defer ts.Close()

	tests := []struct {
		description        string
		url                string
		expectedStatusCode int
		initService        func() validator.Service
	}{
		{
			description:        "404 not found",
			url:                "%s/v2/validate",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description:        "404 not found",
			url:                "%s/validator/vat",
			expectedStatusCode: http.StatusNotFound, // we are testing this in build time
		},
		{
			description:        "417 expectation failed",
			url:                "%s/validator/vat/sd45f768",
			expectedStatusCode: http.StatusExpectationFailed,
		},
		{
			description:        "200 ok",
			url:                "%s/validator/vat/DE156129043",
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "417 expectation failed",
			url:                "%s/validator/vat/De156129043",
			expectedStatusCode: http.StatusExpectationFailed,
		},
		{
			description:        "417 expectation failed",
			url:                "%s/validator/vat/1de156129043",
			expectedStatusCode: http.StatusExpectationFailed,
		},
		{
			description:        "200 ok",
			url:                "%s/validator/vat/DE1561290423",
			expectedStatusCode: http.StatusExpectationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			url := fmt.Sprintf(tt.url, ts.URL)
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			assert.Equal(t, resp.StatusCode, tt.expectedStatusCode)
		})
	}
}
