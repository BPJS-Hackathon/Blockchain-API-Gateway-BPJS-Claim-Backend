package handlers

import (
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(engine *gin.Engine, authService services.AuthService) {
	handler := &AuthHandler{authService: authService}
	// public routes
	gin := engine.Group("/auth")
	gin.POST("/login", handler.Login)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "invalid request",
			"success": false,
			"message": "Login Failed",
		})
		return
	}
	token, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"error":   "unauthorized",
			"success": false,
			"message": "Login Failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"token":   token,
		"success": true,
		"message": "Login Successful",
	})
}
