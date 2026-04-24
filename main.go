package main

import (
	"log"
	"task-manager/db"
	"task-manager/handler"
	"task-manager/repo"
	"task-manager/service"

	"github.com/gin-gonic/gin"
)

func main() {
	gormDSN := "host=localhost user=postgres password=Ernar17042006 dbname=todo port=5432 sslmode=disable"

	migrateDSN := "postgres://postgres:Ernar17042006@localhost:5432/todo?sslmode=disable"

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

	userRepo := repo.NewUserRepository(database)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	authHandler.RegisterRoutes(r)
	todoHandler.RegisterRoutes(r)
	r.Run(":8080")
}
