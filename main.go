package main

import (
	"log"
	"task-manager/db"
	"task-manager/handler"
	"task-manager/service"

	"task-manager/repo"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "host=localhost user=postgres password=Ernar17042006 dbname=todo port=5432 sslmode=disable"

	db, err := db.InitDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	repo := repo.NewTodoRepository(db)
	service := service.NewTodoService(repo)
	handler := handler.NewTodoHandler(service)

	r := gin.Default()
	handler.RegisterRoutes(r)

	r.Run(":8080")
}
