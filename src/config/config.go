package config

import "log"

type AppConfig struct {
	GoEnv string
	InfoLog *log.Logger
	ErrorLog *log.Logger
}