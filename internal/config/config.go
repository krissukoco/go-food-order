package config

import (
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func LoadAndValidate(target interface{}, filenames ...string) error {
	err := godotenv.Overload(filenames...)
	if err != nil {
		log.Printf("error loading env file: %v", err)
	}
	if err := env.Parse(target); err != nil {
		return err
	}
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(target); err != nil {
		return err
	}
	return nil
}
