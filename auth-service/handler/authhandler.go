package handler

import (
	"auth-service/models"
	"auth-service/service"
	"auth-service/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service   *service.AuthService
	jwtSecret string
}

func NewAuthHandler(s *service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{Service: s, jwtSecret: jwtSecret}
}

func (h *AuthHandler) AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.GET("/validate", h.Validate)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	added, err := h.Service.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": added,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Validate(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	claims, err := utils.ValidateToken(token, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	userID := int(claims["user_id"].(float64))

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"user_id": userID,
	})
}
