package healthcheck

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"orange-go/internal/api/healthcheck/model"
)

type API struct {
}

func NewAPI(e *echo.Echo) *API {
	api := &API{}

	e.GET("/live", api.Liveness)
	e.GET("/ready", api.Readiness)

	return api
}

// Liveness godoc
// @Summary      Проверка доступен (liveness)
// @Description  Возвращает статус 200, если сервис доступен.
// @Tags         Healthcheck
// @Accept       json
// @Produce      json
// @Success      200  {object} model.Success "Операция выполнена успешно"
// @Router       /live [get]
func (api *API) Liveness(c echo.Context) error {
	return c.JSON(http.StatusOK, &model.Success{Data: model.SuccessData{
		Message: "OK",
		Code:    http.StatusOK,
	}})
}

func (api *API) Readiness(c echo.Context) error {
	return c.JSON(http.StatusOK, &model.Success{Data: model.SuccessData{
		Message: "OK",
		Code:    http.StatusOK,
	}})
}
