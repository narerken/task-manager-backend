package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type ServerConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode"`
}

type MainServiceConfig struct {
	BaseURL        string `json:"base_url"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

type InternalServiceConfig struct {
	Token string `json:"token"`
}

type Config struct {
	Server          ServerConfig          `json:"server"`
	Database        DatabaseConfig        `json:"database"`
	MainService     MainServiceConfig     `json:"main_service"`
	InternalService InternalServiceConfig `json:"internal_service"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	ApplyDefaults(&cfg)
	applyEnvOverrides(&cfg)
	return &cfg, nil
}

func ApplyDefaults(cfg *Config) {
	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8081
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 10
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 10
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}
	if cfg.MainService.TimeoutSeconds == 0 {
		cfg.MainService.TimeoutSeconds = 10
	}
}

func applyEnvOverrides(cfg *Config) {
	applyStringEnv("SERVER_HOST", &cfg.Server.Host)
	applyIntEnv("SERVER_PORT", &cfg.Server.Port)
	applyIntEnv("SERVER_READ_TIMEOUT", &cfg.Server.ReadTimeout)
	applyIntEnv("SERVER_WRITE_TIMEOUT", &cfg.Server.WriteTimeout)

	applyStringEnv("DATABASE_HOST", &cfg.Database.Host)
	applyIntEnv("DATABASE_PORT", &cfg.Database.Port)
	applyStringEnv("DATABASE_USER", &cfg.Database.User)
	applyStringEnv("DATABASE_PASSWORD", &cfg.Database.Password)
	applyStringEnv("DATABASE_NAME", &cfg.Database.DBName)
	applyStringEnv("DATABASE_SSLMODE", &cfg.Database.SSLMode)

	applyStringEnv("MAIN_SERVICE_BASE_URL", &cfg.MainService.BaseURL)
	applyIntEnv("MAIN_SERVICE_TIMEOUT_SECONDS", &cfg.MainService.TimeoutSeconds)

	applyStringEnv("INTERNAL_SERVICE_TOKEN", &cfg.InternalService.Token)
}

func applyStringEnv(key string, target *string) {
	if value := os.Getenv(key); value != "" {
		*target = value
	}
}

func applyIntEnv(key string, target *int) {
	value := os.Getenv(key)
	if value == "" {
		return
	}

	parsed, err := strconv.Atoi(value)
	if err == nil {
		*target = parsed
	}
}
