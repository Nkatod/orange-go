package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"orange-go/internal/config"
	"orange-go/internal/injection"
	"orange-go/internal/logger"
	tl "orange-go/internal/transactions/logger"
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

	app, err := injection.InitializeApp(e)
	if err != nil {
		log.Fatal(err)
	}
	if err = tl.InitializeTransactionLog(app.TransactionLogger); err != nil {
		log.Fatal(err)
	}

	// ВАЖНО: defer теперь сработает, потому что мы не вызываем os.Exit
	defer func() {
		fmt.Printf("Closing transaction log file\n")
		if err := app.TransactionLogger.Close(); err != nil {
			log.Printf("transaction logger close error: %v", err)
		}
	}()

	// Запускаем сервер в горутине
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server error: %v", err)
		}
	}()

	// Ждем Ctrl+C / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Плавное выключение
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
