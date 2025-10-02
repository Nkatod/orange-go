package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"

	"orange-go/internal/config"
	"orange-go/internal/injection"
	"orange-go/internal/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// инициализируем логер
	if err := logger.Initialize(cfg.LogLevel, "orange-go-app", cfg.Version); err != nil {
		panic(err)
	}

	e := echo.New()

	injection.InitializeHealthCheckAPI(e)
	injection.InitializeKeyValueAPI(e)

	log.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
