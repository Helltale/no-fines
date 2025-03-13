package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	APP_HTTP_PORT string `env:"APP_HTTP_PORT" envDefault:"8080" validate:"required,numeric,gte=1"`
	APP_GRPC_PORT string `env:"APP_GRPC_PORT" envDefault:"50051" validate:"required,numeric,gte=1"`

	DB_CONNECTION_INITIAL_DELAY int `env:"DB_CONNECTION_INITIAL_DELAY" envDefault:"5" validate:"required,gte=0"`
	DB_CONNECTION_MULTIPLIER    int `env:"DB_CONNECTION_MULTIPLIER" envDefault:"2" validate:"required,gte=1"`
	DB_CONNECTION_MAX_DELAY     int `env:"DB_CONNECTION_MAX_DELAY" envDefault:"60" validate:"required,gte=1"`
	DB_CONNECTION_MAX_ATTEMPTS  int `env:"DB_CONNECTION_MAX_ATTEMPTS" envDefault:"5" validate:"required,gte=1"`

	POSTGRES_PORT     string `env:"POSTGRES_PORT,required" validate:"required,numeric,gte=1"`
	POSTGRES_HOST     string `env:"POSTGRES_HOST,required" validate:"required,min=1,max=100"`
	POSTGRES_NAME     string `env:"POSTGRES_NAME,required" validate:"required,min=1,max=100"`
	POSTGRES_USER     string `env:"POSTGRES_USER,required" validate:"required,min=1,max=100"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD,required" validate:"required,min=1,max=100"`
	POSTGRES_SSL_MODE string `env:"POSTGRES_SSL_MODE,required" validate:"required,oneof=enable disable"`

	POSTGRES_MIGRATION_PACKAGE string `env:"POSTGRES_MIGRATION_PACKAGE,required" validate:"required,min=1,max=500"`
}

var Conf Config
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func LoadEnv() (*Config, error) {
	if err := env.Parse(&Conf); err != nil {
		return nil, fmt.Errorf("failed to load the env: %v", err)
	}
	if err := validate.Struct(Conf); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}
	return &Conf, nil
}
