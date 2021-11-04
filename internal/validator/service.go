package validator

import (
	"VatIdValidator/internal/logger"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	VatValidator(vatNum string) (bool, error)
}

type service struct {
	validate *validator.Validate
	logger   *logger.StandardLogger
	config   *Config
	app      AppRepository
}

func CreateService(
	config *Config,
	app AppRepository,
	logger *logger.StandardLogger,
	validator *validator.Validate) Service {
	return &service{
		validate: validator,
		app:      app,
		logger:   logger,
		config:   config,
	}
}
