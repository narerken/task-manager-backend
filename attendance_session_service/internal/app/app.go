package app

import (
	"attendance_session_service/internal/client"
	"attendance_session_service/internal/config"
	"attendance_session_service/internal/database"
	"attendance_session_service/internal/handler"
	"attendance_session_service/internal/repository"
	"attendance_session_service/internal/routes"
	"attendance_session_service/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql db: %v", err)
	}
	defer sqlDB.Close()

	repo := repository.NewAttendanceSessionRepository(db)
	mainClient := client.NewMainServiceClient(
		cfg.MainService.BaseURL,
		cfg.InternalService.Token,
		cfg.MainService.TimeoutSeconds,
	)
	service := service.NewAttendanceSessionService(repo, mainClient)
	handler := handler.NewAttendanceSessionHandler(service)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	routes.RegisterAttendanceSessionRoutes(router, handler, cfg.InternalService.Token)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(router.Run(addr))
}
