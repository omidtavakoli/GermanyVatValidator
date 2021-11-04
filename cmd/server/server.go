package main

import (
	"VatIdValidator/internal/http/rest"
	"VatIdValidator/internal/logger"
	validationService "VatIdValidator/internal/validator"
	"VatIdValidator/pkg/EU_VIES"
	"context"
	"github.com/go-playground/validator/v10"
	"os"
	"sync"
)

type Server struct {
	sync.WaitGroup
	Config      *MainConfig
	RESTHandler *rest.Handler
	Logger      *logger.StandardLogger
}

func NewServer(cfg *MainConfig, logger *logger.StandardLogger) *Server {
	return &Server{
		Config: cfg,
		Logger: logger,
	}
}

// Initialize is responsible for EU_VIES initialization and wrapping required dependencies
func (s *Server) Initialize(ctx context.Context) error {
	v := validator.New()

	appClient, err := EU_VIES.NewClient(&s.Config.EuVies)
	if err != nil {
		return err
	}

	service := validationService.CreateService(&s.Config.Validator, appClient, s.Logger, v)
	handler := rest.CreateHandler(service)
	s.RESTHandler = handler
	return nil
}

// Start starts the application in blocking mode
func (s *Server) Start(ctx context.Context) {
	// Create Router for HTTP Server
	router := SetupRouter(s.RESTHandler, s.Config)

	// Start REST Server in Blocking mode
	s.RESTHandler.Start(ctx, s.Config.Server.Port, router)
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the EU_VIES
func (s *Server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	// Wait for OS signals
	<-quitSignal

	// Kill the API Endpoints first
	s.RESTHandler.Stop()

	close(done)
}
