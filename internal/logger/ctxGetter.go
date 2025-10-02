package logger

import (
	"context"

	"time"

	"orange-go/internal/errors"
)

var serviceName string

// SetServiceName задаёт имя сервиса (вызывается один раз при инициализации).
func SetServiceName(name string) {
	serviceName = name
}

// GetServiceName возвращает имя сервиса.
func GetServiceName() string {
	return serviceName
}

type ContextKey string

// RequestContextKey Ключ для хранения данных запроса в контексте
const RequestContextKey ContextKey = "UserRequestCTXkey"

// RequestDTO Данная структура для хранения общих данных запроса
type RequestDTO struct {
	Method    string
	URL       string
	StartTime time.Time // Время когда поступил запрос
	RequestID string    // Идентификатор запроса (X-Request-Id)
	Cached    bool
}

// GetRequestDataFromCtx Получение данных запроса пользователя из контекста
func GetRequestDataFromCtx(ctx context.Context) (*RequestDTO, error) {
	requestData, ok := ctx.Value(RequestContextKey).(*RequestDTO)
	if !ok {
		return nil, errors.NewCommonError("failed to get request data from context", errors.E1001)
	}
	return requestData, nil
}
