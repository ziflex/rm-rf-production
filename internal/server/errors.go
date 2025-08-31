package server

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ziflex/rm-rf-production/pkg/common"
	"github.com/ziflex/rm-rf-production/pkg/transactions"
)

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewApiError(code, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
	}
}

func NewApiErrorFrom(code string, cause error) *ApiError {
	return &ApiError{
		Code:    code,
		Message: cause.Error(),
	}
}

func errorHandler(err error, c echo.Context) {
	log := zerolog.Ctx(c.Request().Context())

	if errors.Is(err, common.ErrNotFound) {
		c.JSON(404, NewApiErrorFrom("notFound", err))
	} else if errors.Is(err, common.ErrDuplicate) {
		c.JSON(409, NewApiErrorFrom("duplicate", err))
	} else if errors.Is(err, transactions.ErrInvalidOperationType) {
		c.JSON(400, NewApiErrorFrom("invalidOperationType", err))
	} else if he, ok := err.(*echo.HTTPError); ok {
		c.JSON(he.Code, NewApiError("badRequest", he.Message.(string)))
	} else {
		log.Err(err).Msg("internal server error")
		c.JSON(500, echo.Map{"error": "internal server error"})
	}
}
