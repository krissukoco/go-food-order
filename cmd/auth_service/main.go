package main

import (
	"context"
	"log"
	"time"

	"github.com/krissukoco/go-food-order-microservices/internal/app"
	"github.com/krissukoco/go-food-order-microservices/internal/config"
)

type Config struct {
	Environment      string `env:"ENVIRONMENT"  envDetault:"dev"`
	RestPort         uint   `env:"REST_PORT" envDefault:"31000"`
	GrpcPort         uint   `env:"GRPC_PORT" envDefault:"51000"`
	GracefulWaitTime uint   `env:"GRACEFUL_WAIT_TIME" envDefault:"10"`
}

func main() {
	var cfg Config
	if err := config.LoadAndValidate(&cfg); err != nil {
		log.Fatalf("ERROR loading config: %v", err)
	}

	s := app.NewServer(
		cfg.RestPort,
		cfg.GrpcPort,
		time.Duration(cfg.GracefulWaitTime)*time.Second,
		nil,
	)
	if err := s.ListenAndStartGracefully(context.Background()); err != nil {
		panic(err)
	}
}
