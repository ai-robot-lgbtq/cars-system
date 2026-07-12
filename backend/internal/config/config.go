package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv  string
	AppPort int

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	WSBroker       string
	PaymentGateway string
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", 8080)

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "cars")
	v.SetDefault("DB_PASSWORD", "cars_pass")
	v.SetDefault("DB_NAME", "cars_db")
	v.SetDefault("DB_SSLMODE", "disable")

	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("JWT_SECRET", "dev-secret-change-me-in-production")
	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h")

	v.SetDefault("WS_BROKER", "memory")
	v.SetDefault("PAYMENT_GATEWAY", "local")
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	cfg := &Config{
		AppEnv:         v.GetString("APP_ENV"),
		AppPort:        v.GetInt("APP_PORT"),
		DBHost:         v.GetString("DB_HOST"),
		DBPort:         v.GetInt("DB_PORT"),
		DBUser:         v.GetString("DB_USER"),
		DBPassword:     v.GetString("DB_PASSWORD"),
		DBName:         v.GetString("DB_NAME"),
		DBSSLMode:      v.GetString("DB_SSLMODE"),
		RedisHost:      v.GetString("REDIS_HOST"),
		RedisPort:      v.GetInt("REDIS_PORT"),
		RedisPassword:  v.GetString("REDIS_PASSWORD"),
		RedisDB:        v.GetInt("REDIS_DB"),
		JWTSecret:      v.GetString("JWT_SECRET"),
		JWTAccessTTL:   v.GetDuration("JWT_ACCESS_TTL"),
		JWTRefreshTTL:  v.GetDuration("JWT_REFRESH_TTL"),
		WSBroker:       v.GetString("WS_BROKER"),
		PaymentGateway: v.GetString("PAYMENT_GATEWAY"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config invalid: %w", err)
	}

	return cfg, nil
}
