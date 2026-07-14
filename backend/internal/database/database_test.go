package database

import (
	"os"

	"testing"

	"github.com/scutech/cars-system/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect_InvalidHost(t *testing.T) {
	os.Setenv("DB_HOST", "nonexistent.invalid.host")
	os.Setenv("DB_PORT", "1")
	os.Setenv("APP_ENV", "test")

	cfg, err := config.Load()
	require.NoError(t, err)

	db, err := Connect(cfg)
	assert.Error(t, err, "expected connnection to fail with invalid host")
	assert.Nil(t, db)

	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("APP_ENV")
}
