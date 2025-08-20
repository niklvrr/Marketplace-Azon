package errs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	NotFoundError      = errors.New("not found")
	NotAuthorizedError = errors.New("not authorized")
	ForbiddenError     = errors.New("forbidden")
	ValidationError    = errors.New("validation error")
)

func RespondError(ctx *gin.Context, status int, code string, message string) {
	ctx.JSON(status, gin.H{
		"data":  nil,
		"error": gin.H{"code": code, "message": message},
	})
}

func RespondServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, NotFoundError):
		RespondError(ctx, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, NotAuthorizedError):
		RespondError(ctx, http.StatusUnauthorized, "unauthorized", err.Error())
	case errors.Is(err, ForbiddenError):
		RespondError(ctx, http.StatusForbidden, "forbidden", err.Error())
	case errors.Is(err, ValidationError):
		RespondError(ctx, http.StatusBadRequest, "validation_error", err.Error())
	default:
		RespondError(ctx, http.StatusInternalServerError, "internal_error", err.Error())
	}
}
