package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *Logger
var version = "unknown"

// Logger является оберткой над logrus.Logger и содержит дополнительное поле "component".
type Logger struct {
	Logrus    *logrus.Logger
	component string
}

// SetVersion устанавливает глобальную версию для логов.
func SetVersion(v string) {
	version = v
}

// GetVersion получает глобальную версию для логов.
func GetVersion() string {
	return version
}

// Info записывает информационное сообщение с полями и дополнительными данными из контекста.
func (l *Logger) Info(ctx context.Context, msg string, fields ...logrus.Fields) {
	// Если поля не переданы, используем пустой набор полей.
	var f logrus.Fields
	if len(fields) > 0 {
		f = fields[0]
	} else {
		f = logrus.Fields{}
	}
	LogWithUserData(ctx).WithFields(f).Info(msg)
}

// Error записывает сообщение об ошибке с полями и дополнительными данными из контекста.
func (l *Logger) Error(ctx context.Context, msg string, fields ...logrus.Fields) {
	// Если поля не переданы, используем пустой набор полей.
	var f logrus.Fields
	if len(fields) > 0 {
		f = fields[0]
	} else {
		f = logrus.Fields{}
	}
	LogWithUserData(ctx).WithFields(f).Error(msg)
}

// Initialize инициализирует глобальный синглтон-логгер с указанным уровнем и названием компонента.
func Initialize(level, component, serviceVersion string) error {
	// Определяем уровень логирования по строковому значению.
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	// Создаем новый логгер.
	log := logrus.New()
	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)

	// Настраиваем JSON форматтер и мэппинг ключей.
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000-0700",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",    // ключ для времени
			logrus.FieldKeyLevel: "level",   // ключ для уровня логирования
			logrus.FieldKeyMsg:   "message", // ключ для сообщения
		},
	})

	// Сохраняем глобальный логгер.
	logger = &Logger{
		Logrus:    log,
		component: component,
	}

	return nil
}

// addCallerField добавляет в logrus.Entry информацию о вызывающем месте, используя заданное значение skip.
func addCallerField(entry *logrus.Entry, skip int) *logrus.Entry {
	// skip = 0: runtime.Callers, 1: addCallerField, 2: функция, которая вызвала addCallerField.
	if _, file, line, ok := runtime.Caller(skip); ok {
		return entry.WithField("caller", fmt.Sprintf("%s:%d", file, line))
	}
	return entry.WithField("caller", "unknown")
}

// LogWithUserData добавляет данные пользователя из контекста к логированию и возвращает logrus.Entry.
func LogWithUserData(ctx context.Context) *logrus.Entry {
	// Начинаем с глобального логгера и добавляем поле component.
	entry := Get().Logrus.WithField("component", Get().component)

	// Добавляем информацию о вызове – здесь skip равен 3, так как:
	// 0: runtime.Caller внутри addCallerField,
	// 1: addCallerField,
	// 2: LogWithUserData,
	// 3: та функция, которая вызвала LogWithUserData.
	entry = addCallerField(entry, 3)

	// Извлекаем данные запроса из контекста.
	method := ""
	url := ""
	requestID := ""
	cutoff := 0.0

	entry = entry.WithFields(logrus.Fields{
		"version": version,
	})

	if requestData, err := GetRequestDataFromCtx(ctx); err == nil {
		method = requestData.Method
		url = requestData.URL
		requestID = requestData.RequestID
		cutoff = time.Since(requestData.StartTime).Seconds() // преобразуем в секунды с плавающей точкой
		entry = entry.WithFields(logrus.Fields{
			"request-id":   requestID,
			"cutoff":       cutoff,
			"handler-name": method + " " + url,
		})
	}

	return entry
}

// NopLogger возвращает логгер, который не выполняет никаких действий (выводим в os.Discard).
func NopLogger() *Logger {
	nopLog := logrus.New()
	return &Logger{Logrus: nopLog}
}

// Get возвращает текущий глобальный логгер, или NopLogger, если он не инициализирован.
func Get() *Logger {
	if logger == nil {
		return NopLogger()
	}
	return logger
}
