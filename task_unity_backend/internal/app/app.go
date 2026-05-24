package app

import (
	"fmt"
	"log"
	"task-manager/internal/auth"
	"task-manager/internal/client"
	"task-manager/internal/config"
	"task-manager/internal/database"
	"task-manager/internal/handler"
	"task-manager/internal/middleware"
	"task-manager/internal/repository"
	"task-manager/internal/routes"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load data from .env: %v", err)
	}

	err = config.ValidateConfig(cfg)
	if err != nil {
		log.Fatalf("invalid token: %v", err)
	}

	pool, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	sqlDB, err := pool.DB()
	if err != nil {
		log.Fatalf("failed to get sql db: %v", err)
	}

	defer sqlDB.Close()

	fmt.Println("Success", pool)

	jwtSecret := auth.JWTSecret{Secret: []byte(cfg.JWTSecret)}
	userR := repository.NewUserRepository(pool)
	userS := service.NewUserService(userR, &jwtSecret)
	userH := handler.NewUserHandler(userS)

	taskR := repository.NewTaskRepo(pool)
	taskS := service.NewTaskService(taskR, userR)
	taskH := handler.NewTaskHandler(taskS)

	departmentR := repository.NewDepartmentRepository(pool)
	departmentS := service.NewDepartmentService(departmentR)
	departmentH := handler.NewDepartmentHandler(departmentS)

	attendanceR := repository.NewAttendanceRepository(pool)
	attendanceS := service.NewAttendanceService(attendanceR, userS)
	attendanceH := handler.NewAttendanceHandler(attendanceS)

	attendanceSessionClient := client.NewAttendanceSessionClient(
		cfg.AttendanceSessionService.BaseURL,
		cfg.InternalService.Token,
		cfg.AttendanceSessionService.TimeoutSeconds,
	)
	attendanceSessionH := handler.NewAttendanceSessionHandler(attendanceSessionClient)

	internalDataH := handler.NewInternalDataHandler(userR, departmentR)

	commentR := repository.NewCommentRepository(pool)
	commentS := service.NewCommentService(commentR, taskS)
	commentH := handler.NewCommentHandler(commentS)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS())

	routes.UserRoutes(userH, router, &jwtSecret)
	routes.TaskRoutes(taskH, router, &jwtSecret)
	routes.DepartmentRoutes(departmentH, router, &jwtSecret)
	routes.AttendanceRoutes(attendanceH, router, &jwtSecret)
	routes.AttendanceSessionRoutes(attendanceSessionH, router, &jwtSecret)
	routes.CommentRoutes(commentH, router, &jwtSecret)
	routes.InternalServiceRoutes(internalDataH, router, cfg.InternalService.Token)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(router.Run(addr))
}
