//go:build wireinject
// +build wireinject

package injection

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"

	healthcheck "orange-go/internal/api/healthcheck"
	key_value_api "orange-go/internal/api/key"
	"orange-go/internal/storage"
	"orange-go/internal/storage/memory"
)

func InitializeHealthCheckAPI(e *echo.Echo) *healthcheck.API {
	wire.Build(
		healthcheck.NewAPI,
	)

	return &healthcheck.API{}
}

func InitializeKeyValueAPI(e *echo.Echo) *key_value_api.API {
	wire.Build(
		// провайдеры зависимостей key.NewAPI:
		memory.NewMemoryStorage,
		wire.Bind(new(storage.IMemoryStorage), new(*memory.Storage)), // IMemoryStorage <- *MemoryStorage

		key_value_api.NewAPI, // (*echo.Echo, storage.IMemoryStorage) *API
	)
	return &key_value_api.API{}
}
