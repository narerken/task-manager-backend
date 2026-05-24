package routes

import (
	"attendance_session_service/internal/handler"
	"attendance_session_service/internal/httpx"
	"attendance_session_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAttendanceSessionRoutes(router gin.IRoutes, h *handler.AttendanceSessionHandler, serviceToken string) {
	router.POST("/attendance/sessions", httpx.WrapHandler(middleware.ServiceAuthMiddleware(serviceToken, h.CreateSession)))
	router.GET("/attendance/sessions", httpx.WrapHandler(middleware.ServiceAuthMiddleware(serviceToken, h.ListSessions)))
	router.GET("/attendance/sessions/:id", httpx.WrapHandler(middleware.ServiceAuthMiddleware(serviceToken, h.GetSession)))
	router.PATCH("/attendance/sessions/:id/entries", httpx.WrapHandler(middleware.ServiceAuthMiddleware(serviceToken, h.BulkUpsertEntries)))
	router.PATCH("/attendance/sessions/:id/publish", httpx.WrapHandler(middleware.ServiceAuthMiddleware(serviceToken, h.PublishSession)))
}
