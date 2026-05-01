package main

import (
	"fmt"
	"log"
	"todo-service/client"
	"todo-service/config"
	"todo-service/db"
	"todo-service/handler"
	"todo-service/repo"
	"todo-service/service"

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

	todoRepo := repo.NewTodoRepository(database)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	r := gin.Default()

	authClient := client.NewAuthClient(cfg.AuthServiceURL)
	todoHandler.TodoRoutes(r, authClient)
	err = r.Run(fmt.Sprintf(":%d", cfg.ServerCfg.Port))
	if err != nil {
		log.Fatal("server run failed:", err)
	}
}
