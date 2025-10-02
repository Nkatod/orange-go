package key

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"orange-go/internal/model"
	"orange-go/internal/storage"
	"orange-go/internal/transactions"
)

type API struct {
	memoryStorage storage.IMemoryStorage
	tlog          transactions.ITransactionLogger
}

func NewAPI(e *echo.Echo, memoryStorage storage.IMemoryStorage, tlog transactions.ITransactionLogger) *API {
	api := &API{memoryStorage: memoryStorage, tlog: tlog}

	g := e.Group("/v1/key")
	g.PUT("/:key", api.KeyValuePutHandler)
	g.GET("/:key", api.KeyValueGetHandler)

	return api
}

// KeyValuePutHandler godoc
// @Summary      Создание/обновление key-value
// @Description  Сохраняет сырое значение (body) по ключу в памяти
// @Tags         Storage
// @Accept       */*
// @Produce      json
// @Param        key   path      string  true  "Key"
// @Success      201   {object}  model.KeyValuePutResponse
// @Failure      400   {object}  model.KeyValuePutResponse
// @Failure      500   {object}  model.KeyValuePutResponse
// @Router       /v1/key/{key} [put]
func (api *API) KeyValuePutHandler(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, KeyValuePutResponse{
			Error: model.NewErr("bad_request", "empty key", ""),
		})
	}

	defer func() { _ = c.Request().Body.Close() }()
	valueBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, KeyValuePutResponse{
			Error: model.NewErr("read_body_failed", "failed to read request body", err.Error()),
		})
	}

	if err := api.memoryStorage.Put(key, string(valueBytes)); err != nil {
		return c.JSON(http.StatusInternalServerError, KeyValuePutResponse{
			Error: model.NewErr("storage_put_failed", "failed to store value", err.Error()),
		})
	}

	api.tlog.WritePut(key, string(valueBytes))

	return c.JSON(http.StatusCreated, KeyValuePutResponse{
		Data: &KeyValuePutData{
			Key:    key,
			Status: "created",
		},
	})
}

// KeyValueGetHandler godoc
// @Summary      Создание/обновление key-value
// @Description  Сохраняет сырое значение (body) по ключу в памяти
// @Tags         Storage
// @Accept       */*
// @Produce      json
// @Param        key   path      string  true  "Key"
// @Success      201   {object}  model.KeyValueGetResponse
// @Failure      400   {object}  model.KeyValueGetResponse
// @Failure      500   {object}  model.KeyValueGetResponse
// @Router       /v1/key/{key} [get]
func (api *API) KeyValueGetHandler(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, KeyValueGetResponse{
			Error: model.NewErr("bad_request", "empty key", ""),
		})
	}

	value, err := api.memoryStorage.Get(key)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, KeyValueGetResponse{
			Error: model.NewErr("storage_put_failed", "failed to store value", err.Error()),
		})
	}

	return c.JSON(http.StatusCreated, KeyValueGetResponse{
		Data: value,
	})
}
