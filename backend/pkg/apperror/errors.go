package apperror

import "errors"

// общие ошибки приложения, используемые на уровне бизнес-логики
// чтобы абстрагировать слой domain от конкретной реализации хранилища
// например, когда хранилище возврощает специфичную ошибку (clean architecture)
var (
	ErrNotFound           = errors.New("ресурс не найден")
	ErrAlreadyExists      = errors.New("ресурс уже существует")
	ErrInvalidCredentials = errors.New("неверные учетные данные")
	ErrUnauthorized       = errors.New("не авторизован")
	ErrForbidden          = errors.New("доступ запрещен")
	ErrInvalidInput       = errors.New("некорректные входные данные")
	ErrInternal           = errors.New("внутрення ошибка сервера")
)

// для унифицированного ответа клиенту при ошибках
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// реализация интерфейса error для AppError
func (e *AppError) Error() string {
	return e.Message
}

// создает новую ошибку приложения с HTTP-кодом и сообщением
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
