package main

import (
	"VatIdValidator/internal/logger"
	"VatIdValidator/internal/validator"
	"VatIdValidator/pkg/EU_VIES"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	User        string `yaml:"USER"`
	Pass        string `yaml:"PASS"`
	AuthEnabled bool   `yaml:"AUTH_ENABLED"`
	Port        int32  `yaml:"PORT"`
}

type MainConfig struct {
	Logger logger.Config    `yaml:"LOGGER"`
	Proxy  validator.Config `yaml:"VAT_ID_VALIDATOR"`
	Server ServerConfig     `yaml:"SERVER"`
	App    EU_VIES.Config   `yaml:"APP"`
}

// LoadConfig loads configs form provided yaml file or overrides it with env variables
func LoadConfig(filePath string) (*MainConfig, error) {
	cfg := MainConfig{}
	if filePath != "" {
		err := readFile(&cfg, filePath)
		if err != nil {
			return nil, err
		}
	}
	err := readEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func readFile(cfg *MainConfig, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func readEnv(cfg *MainConfig) error {
	return envconfig.Process("", cfg)
}
