package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

type AttendanceSessionServiceConfig struct {
	BaseURL        string `json:"base_url"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

type InternalServiceConfig struct {
	Token string `json:"token"`
}

type Config struct {
	Server                   ServerConfig                   `json:"server"`
	Database                 DatabaseConfig                 `json:"database"`
	AttendanceSessionService AttendanceSessionServiceConfig `json:"attendance_session_service"`
	InternalService          InternalServiceConfig          `json:"internal_service"`
	JWTSecret                string                         `json:"jwt_secret"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	cfg := Config{}

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	ApplyDefaults(&cfg)
	applyEnvOverrides(&cfg)
	return &cfg, nil
}

func ValidateConfig(cfg *Config) error {
	var errors []string

	if cfg.Server.Host == "" {
		errors = append(errors, "server_host обязателен")
	}

	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		errors = append(errors, "недопустимое значение для server_port")
	}

	if cfg.Database.Host == "" {
		errors = append(errors, "database_host обязателен")
	}

	if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
		errors = append(errors, "недопустимое значение для database_port")
	}

	if cfg.Database.User == "" {
		errors = append(errors, "database_user обязателен")
	}

	if len(errors) > 0 {
		return fmt.Errorf("ошибки валидации: %s", strings.Join(errors, "; "))
	}
	return nil
}

func ApplyDefaults(cfg *Config) {
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30
	}

	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30
	}

	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}

	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disabled"
	}

	if cfg.AttendanceSessionService.TimeoutSeconds == 0 {
		cfg.AttendanceSessionService.TimeoutSeconds = 10
	}
}

func NewDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Database: DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			SSLMode: "disable",
		},
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

	applyStringEnv("ATTENDANCE_SESSION_SERVICE_BASE_URL", &cfg.AttendanceSessionService.BaseURL)
	applyIntEnv("ATTENDANCE_SESSION_SERVICE_TIMEOUT_SECONDS", &cfg.AttendanceSessionService.TimeoutSeconds)

	applyStringEnv("INTERNAL_SERVICE_TOKEN", &cfg.InternalService.Token)
	applyStringEnv("JWT_SECRET", &cfg.JWTSecret)
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
