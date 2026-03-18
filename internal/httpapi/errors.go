package httpapi

import (
	"net/http"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/observability"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteError(c *gin.Context, err error) {
	status, payload := buildErrorResponse(err)
	logError(c, err, status)
	c.JSON(status, payload)
}

func AbortWithError(c *gin.Context, err error) {
	status, payload := buildErrorResponse(err)
	logError(c, err, status)
	c.AbortWithStatusJSON(status, payload)
}

func buildErrorResponse(err error) (int, ErrorResponse) {
	appErr, ok := apperrors.As(err)
	if !ok {
		return http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "internal_error",
				Message: "internal server error",
			},
		}
	}

	return statusCode(appErr.Kind), ErrorResponse{
		Error: ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	}
}

func statusCode(kind apperrors.Kind) int {
	switch kind {
	case apperrors.KindValidation:
		return http.StatusBadRequest
	case apperrors.KindUnauthorized:
		return http.StatusUnauthorized
	case apperrors.KindForbidden:
		return http.StatusForbidden
	case apperrors.KindNotFound:
		return http.StatusNotFound
	case apperrors.KindConflict:
		return http.StatusConflict
	case apperrors.KindUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func logError(c *gin.Context, err error, status int) {
	logger := observability.LoggerFromGinContext(c)
	args := []any{"status", status}

	appErr, ok := apperrors.As(err)
	if ok {
		args = append(
			args,
			"error_code", appErr.Code,
			"error_kind", appErr.Kind,
			"error_message", appErr.Message,
		)
		if appErr.Err != nil {
			args = append(args, "error_cause", appErr.Err.Error())
		}
	} else if err != nil {
		args = append(
			args,
			"error_code", "internal_error",
			"error_kind", apperrors.KindInternal,
			"error_message", "internal server error",
			"error_cause", err.Error(),
		)
	}

	if status >= http.StatusInternalServerError {
		logger.Error("request failed", args...)
		return
	}

	logger.Warn("request failed", args...)
}
