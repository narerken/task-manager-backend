package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerCfg      ServerConfig
	DBCfg          DBConfig
	JWTSecret      string
	AuthServiceURL string
}

type ServerConfig struct {
	Port int
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerCfg: ServerConfig{
			Port: serverPort,
		},
		DBCfg: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AuthServiceURL: os.Getenv("AUTH_SERVICE_URL"),
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.ServerCfg.Port == 0 {
		return errors.New("server port is required")
	}

	if c.DBConfigEmpty() {
		return errors.New("db config is incomplete")
	}

	if c.JWTSecret == "" {
		return errors.New("jwt secret is required")
	}

	return nil
}

func (c *Config) DBConfigEmpty() bool {
	db := c.DBCfg

	return db.Host == "" ||
		db.Port == 0 ||
		db.User == "" ||
		db.Password == "" ||
		db.Name == "" ||
		db.SSLMode == ""
}
