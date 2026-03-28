package handler

import (
	"errors"
	"net/http"

	"github.com/alligatorO15/taskMind/backend/pkg/apperror"
	"github.com/gin-gonic/gin"
)

// handlerError - преоьразует ошибки приложения в HTTP-ответы
func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	msg := err.Error()
	code := http.StatusInternalServerError // 500

	switch {
	case errors.Is(err, apperror.ErrAlreadyExists):
		code = http.StatusConflict // 409
	case errors.Is(err, apperror.ErrInvalidCredentials), errors.Is(err, apperror.ErrUnauthorized):
		code = http.StatusUnauthorized // 401
	case errors.Is(err, apperror.ErrNotFound):
		code = http.StatusNotFound // 404
	case errors.Is(err, apperror.ErrInvalidInput):
		code = http.StatusBadRequest // 400
	case errors.Is(err, apperror.ErrForbidden):
		code = http.StatusForbidden // 403
	}

	c.JSON(code, gin.H{"error": msg})
}
