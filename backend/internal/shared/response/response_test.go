package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apperrors "github.com/scutech/cars-system/backend/internal/shared/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ok", func(c *gin.Context) { OK(c, gin.H{"foo": "bar"}) })
	r.GET("/fail", func(c *gin.Context) { Fail(c, apperrors.CodeCarNotFound, "car not found") })
	return r
}

func TestOK(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ok", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.EqualValues(t, 0, body["code"])
	assert.Equal(t, "ok", body["message"])
	assert.NotNil(t, body["data"])
}
func TestFail(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/fail", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code) // 业务错误仍用 200，code 字段标识

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.EqualValues(t, 31001, body["code"])
	assert.Equal(t, "car not found", body["message"])
	assert.Nil(t, body["data"])
}
