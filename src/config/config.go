package config

import (
	"log"

	"github.com/go-playground/validator/v10"
)

type AppConfig struct {
	GoEnv string
	InfoLog *log.Logger
	ErrorLog *log.Logger
	Validate *validator.Validate
}