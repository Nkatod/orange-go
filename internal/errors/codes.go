package errors

type Code uint32

const (
	OK                 Code = iota // OK: http system error map
	Canceled                       // Canceled: http system error map
	Unknown                        // Unknown: http system error map
	InvalidArgument                // InvalidArgument: http system error map
	DeadlineExceeded               // DeadlineExceeded: http system error map
	NotFound                       // NotFound: http system error map
	AlreadyExists                  // AlreadyExists: http system error map
	PermissionDenied               // PermissionDenied: http system error map
	ResourceExhausted              // ResourceExhausted: http system error map
	FailedPrecondition             // FailedPrecondition: http system error map
	Aborted                        // Aborted: http system error map
	OutOfRange                     // OutOfRange: http system error map
	Unimplemented                  // Unimplemented: http system error map
	Internal                       // Internal: http system error map
	Unavailable                    // Unavailable: http system error map
	DataLoss                       // DataLoss: http system error map
	Unauthenticated                // Unauthenticated: http system error map

	// Общие ошибки без бизнес логики
	E1001 = 1001 // E1001 Ошибка получения общих данных из контекста
	E1002 = 1002 // E1002 Отсутствует заголовок Authorization
	E1003 = 1003 // E1003 Неправильный формат токена
	E1004 = 1004 // E1004 Не удалось декодировать claims из токена
	E1005 = 1005 // E1005 Ошибка при разборе JSON
	E1006 = 1006 // E1006 Ошибка при конвертации локальной структуры в UserAuthDTO
	E1007 = 1007 // E1007 Отсутствует заголовок x-capabilities
	E1008 = 1008 // E1008 Ошибка получения авторизации из хранилища
	E1009 = 1009 // E1009 Ошибка преобразование из строки в целое число
	E1010 = 1010 // E1010 Ошибка метрик Prometheus
	E1011 = 1011 // E1011 Ошибка проверки x-capabilities
	E1012 = 1012 // E1012 CheckStringHasOnlyDigits error
	E1013 = 1013 // E1013 json.Unmarshal Error
	E1014 = 1014 // E1014 http transport error
	E1015 = 1015 // E1015 http client do error
	E1016 = 1016 // E1016 http IO.ReadAll error
	E1017 = 1017 // E1017 Mock Error
	E1018 = 1018 // E1018 Cache Key not found
	E1019 = 1019 // E1019 Cache invalid pattern
	E1020 = 1020 // E1020 json.MarshalIndent Error
	E1021 = 1021 // E1021 Неизвестный tenant
	E1022 = 1022 // E1022 Ошибка загрузки временной зоны
	E1023 = 1023 // E1023 Database common error
	E1024 = 1024 // E1023 Database common error

	// Ошибки REST клиента
	E1400 = 1400 // E1400 400 ErrBadRequest
	E1401 = 1401 // E1401 401 ErrUnauthorized
	E1403 = 1403 // E1403 403 ErrForbidden
	E1404 = 1404 // E1404 404 ErrNotFound
	E1422 = 1422 // E1422 422 ErrUnprocessableEntity
	E1450 = 1450 // E1450 500 ErrInternalServerError
	E1452 = 1452 // E1450 502 ErrBadGateway
)
