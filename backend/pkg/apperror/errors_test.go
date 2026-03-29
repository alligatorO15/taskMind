package apperror

import (
	"errors"
	"testing"
)

// TestAppError_Error — проверяет реализацию интерфейса error для AppError
func TestAppError_Error(t *testing.T) {
	appErr := NewAppError(400, "тестовая ошибка")

	if appErr.Error() != "тестовая ошибка" {
		t.Errorf("Неверное сообщение ошибки: '%s'", appErr.Error())
	}
	if appErr.Code != 400 {
		t.Errorf("Неверный код ошибки: %d", appErr.Code)
	}
}

// TestSentinelErrors — проверяет, что sentinel-ошибки корректно сравниваются через errors.Is
func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
		{"ErrAlreadyExists", ErrAlreadyExists},
		{"ErrInvalidCredentials", ErrInvalidCredentials},
		{"ErrUnauthorized", ErrUnauthorized},
		{"ErrForbidden", ErrForbidden},
		{"ErrInvalidInput", ErrInvalidInput},
		{"ErrInternal", ErrInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.err, tt.err) {
				t.Errorf("errors.Is не работает для %s", tt.name)
			}

			if errors.Is(tt.err, errors.New("другая ошибка")) {
				t.Errorf("errors.Is неверно сравнивает %s с другой ошибкой", tt.name)
			}
		})
	}
}
