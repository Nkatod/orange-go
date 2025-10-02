//go:build wireinject
// +build wireinject

package injection

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"

	healthcheck "orange-go/internal/api/healthcheck"
	keyapi "orange-go/internal/api/key"

	"orange-go/internal/storage/memory"
	"orange-go/internal/transactions"
	tl "orange-go/internal/transactions/logger"
)

// Создаем структуру для группировки всех компонентов
type App struct {
	HealthcheckAPI    *healthcheck.API
	KeyAPI            *keyapi.API
	TransactionLogger transactions.ITransactionLogger
}

var storageSet = wire.NewSet(
	memory.NewMemoryStorage, // Убираем wire.Bind, если NewMemoryStorage уже возвращает нужный тип
)

var txLogSet = wire.NewSet(
	wire.Value("transaction.log"),
	tl.NewFileTransactionLogger,
)

// Провайдер для создания App структуры
func NewApp(
	healthAPI *healthcheck.API,
	keyAPI *keyapi.API,
	txLogger transactions.ITransactionLogger,
) *App {
	return &App{
		HealthcheckAPI:    healthAPI,
		KeyAPI:            keyAPI,
		TransactionLogger: txLogger,
	}
}

var appSet = wire.NewSet(
	NewApp,
)

func InitializeApp(e *echo.Echo) (*App, error) {
	wire.Build(
		storageSet,
		txLogSet,
		appSet,
		healthcheck.NewAPI,
		keyapi.NewAPI,
	)
	return nil, nil
}
