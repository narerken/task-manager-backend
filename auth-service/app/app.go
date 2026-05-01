package app

import (
	"auth-service/config"
	"auth-service/db"
	"auth-service/handler"
	"auth-service/repo"
	"auth-service/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("config load failed:", err)
	}

	gormDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.DBCfg.Host,
		cfg.DBCfg.User,
		cfg.DBCfg.Password,
		cfg.DBCfg.Name,
		cfg.DBCfg.Port,
		cfg.DBCfg.SSLMode,
	)

	migrateDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBCfg.User,
		cfg.DBCfg.Password,
		cfg.DBCfg.Host,
		cfg.DBCfg.Port,
		cfg.DBCfg.Name,
		cfg.DBCfg.SSLMode,
	)

	if err := db.RunMigrations(migrateDSN, "file://migrations"); err != nil {
		log.Fatal("migrations failed:", err)
	}

	database, err := db.InitDB(gormDSN)
	if err != nil {
		log.Fatal("db init failed:", err)
	}

	userRepo := repo.NewUserRepository(database)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService, cfg.JWTSecret)

	r := gin.Default()
	authHandler.AuthRoutes(r)
	err = r.Run(fmt.Sprintf(":%d", cfg.ServerCfg.Port))
	if err != nil {
		log.Fatal("server run failed:", err)
	}
}
