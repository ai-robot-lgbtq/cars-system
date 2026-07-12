package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	apperrors "github.com/scutech/cars-system/backend/internal/shared/errors"
	"github.com/scutech/cars-system/backend/internal/shared/response"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", ctx.Request.URL.Path),
					zap.ByteString("stack", debug.Stack()),
				)
				response.Fail(ctx, apperrors.CodeSystemError, "internal server error")
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}
