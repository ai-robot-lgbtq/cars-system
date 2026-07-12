package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/scutech/cars-system/backend/internal/shared/errors"
)

type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Code: apperrors.CodeOK, Message: "ok", Data: data})
}

func Fail(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(http.StatusOK, Envelope{Code: code, Message: message})
}
