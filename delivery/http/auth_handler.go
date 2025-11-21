package handlers

import (
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/config"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/services"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(engine *gin.Engine, authService services.AuthService, jwtManager *utils.JWTManager) {
	handler := &AuthHandler{authService: authService}
	// public routes
	gin := engine.Group("/auth")
	gin.POST("/login", handler.Login)
	forME := engine.Group("/me")
	forME.Use(config.AuthMiddleware(jwtManager))
	forME.GET("", handler.Me)

}

func (h *AuthHandler) Me(c *gin.Context) {
	userHitterID, isVool := c.Get("userID")
	if !isVool {
		c.JSON(500, gin.H{
			"error":   "Unathorized",
			"success": false,
			"message": "Fetch Failed",
		})
		return
	}

	datas, err := h.authService.Me(c.Request.Context(), userHitterID.(string))
	if err != nil {
		c.JSON(401, gin.H{
			"error":   "unauthorized",
			"success": false,
			"message": "Me Checker Failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"data":    datas,
		"success": true,
		"message": "Me",
	})

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
