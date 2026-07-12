package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, 8080, cfg.AppPort)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, 5432, cfg.DBPort)
	assert.Equal(t, "memory", cfg.WSBroker)
	assert.Equal(t, "local", cfg.PaymentGateway)
	assert.Equal(t, 15*time.Minute, cfg.JWTAccessTTL)
	assert.Equal(t, 7*24*time.Hour, cfg.JWTRefreshTTL)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("JWT_SECRET", "supersecret")
	defer clearEnv(t)

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, 9090, cfg.AppPort)
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, "supersecret", cfg.JWTSecret)
}

func clearEnv(t *testing.T) {
	t.Helper()
	vars := []string{
		"APP_ENV", "APP_PORT", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD",
		"DB_NAME", "DB_SSLMODE", "REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD",
		"REDIS_DB", "JWT_SECRET", "JWT_ACCESS_TTL", "JWT_REFRESH_TTL",
		"WS_BROKER", "PAYMENT_GATEWAY",
	}
	for _, v := range vars {
		os.Unsetenv(v)
	}
}
