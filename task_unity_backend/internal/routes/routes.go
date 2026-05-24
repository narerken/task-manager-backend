package routes

import (
	"task-manager/internal/auth"
	"task-manager/internal/handler"
	"task-manager/internal/httpx"
	"task-manager/internal/middleware"
	"task-manager/internal/models"

	"github.com/gin-gonic/gin"
)

func UserRoutes(userH *handler.UserHandler, router gin.IRoutes, jwtSecret *auth.JWTSecret) {
	router.PATCH("/users/update/profile/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, userH.UpdateUserProfile)))
	router.PATCH("/users/update/password/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, userH.UpdateUserPassword)))
	router.DELETE("/users/delete/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(userH.DeleteUser, models.Role{Name: "admin"}))))

	router.POST("/users/register", httpx.WrapHandler(userH.Register))
	router.GET("/users", httpx.WrapHandler(userH.GetAllUsers))
	router.GET("/users/:id", httpx.WrapHandler(userH.GetUserById))
	router.POST("/users/login", httpx.WrapHandler(userH.Login))

}

func TaskRoutes(taskH *handler.TaskHandler, router gin.IRoutes, jwtSecret *auth.JWTSecret) {
	router.POST("/tasks/create", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(taskH.AddTask, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))
	router.GET("/tasks/assignee", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, taskH.GetAllTasksByAssigneeId)))
	router.PATCH("/tasks/update/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(taskH.UpdateTask, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))
	router.DELETE("/tasks/delete/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(taskH.DeleteTask, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.GET("/tasks/:id", httpx.WrapHandler(taskH.GetTaskById))
	router.GET("/tasks", httpx.WrapHandler(taskH.GetAllTasks))
	router.GET("/tasks/", httpx.WrapHandler(taskH.GetAllTasks))
}

func DepartmentRoutes(departmentH *handler.DepartmentHandler, router gin.IRoutes, jwtSecret *auth.JWTSecret) {
	router.POST("/departments/create", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(departmentH.AddDepartment, models.Role{Name: "admin"}))))
	router.GET("/departments", httpx.WrapHandler(departmentH.GetAllDepartments))
	router.GET("/departments/", httpx.WrapHandler(departmentH.GetAllDepartments))
}

func AttendanceRoutes(attendanceH *handler.AttendanceHandler, router gin.IRoutes, jwtSecret *auth.JWTSecret) {
	router.POST("/attendances/create", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceH.AddAttendance, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.GET("/attendances", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceH.GetAllAttendances, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.PATCH("/attendances/update/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceH.UpdateAttendance, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.DELETE("/attendances/delete/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceH.DeleteAttendance, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))
}

func AttendanceSessionRoutes(attendanceSessionH *handler.AttendanceSessionHandler, router gin.IRoutes, jwtSecret *auth.JWTSecret) {
	router.POST("/attendance/sessions", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceSessionH.CreateSession, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.GET("/attendance/sessions", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceSessionH.ListSessions, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.GET("/attendance/sessions/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceSessionH.GetSession, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.PATCH("/attendance/sessions/:id/entries", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceSessionH.BulkUpsertEntries, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.PATCH("/attendance/sessions/:id/publish", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceSessionH.PublishSession, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))
}

func InternalServiceRoutes(internalDataH *handler.InternalDataHandler, router gin.IRoutes, serviceToken string) {
	router.GET("/internal/users/:id", middleware.InternalServiceMiddleware(serviceToken), httpx.WrapHandler(internalDataH.GetUserById))
	router.GET("/internal/departments/:id", middleware.InternalServiceMiddleware(serviceToken), httpx.WrapHandler(internalDataH.GetDepartmentById))
	router.GET("/internal/departments/:id/users", middleware.InternalServiceMiddleware(serviceToken), httpx.WrapHandler(internalDataH.GetUsersByDepartmentId))
}

func CommentRoutes(commentH *handler.CommentHandler, router gin.IRoutes, jwtSecret *auth.JWTSecret) {
	router.POST("/comments/create", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(commentH.AddComment, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.GET("/comments", httpx.WrapHandler(commentH.GetAllComments))

	router.PUT("/comments/update/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(commentH.UpdateComment, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	router.DELETE("/comments/delete/:id", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(commentH.DeleteComment, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))
}
