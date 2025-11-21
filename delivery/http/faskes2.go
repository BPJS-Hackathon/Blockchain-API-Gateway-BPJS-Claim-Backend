package handlers

import (
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/config"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/gin-gonic/gin"
)

type Faskes2Handler struct {
	faskes2Service models.Faskes2Service
}

func NewFaskes2Handler(engine *gin.Engine, faskes2Service models.Faskes2Service, jwtManager *utils.JWTManager) {
	handler := &Faskes2Handler{faskes2Service: faskes2Service}
	// public routes
	gin := engine.Group("/faskes2")
	gin.Use(config.AuthMiddleware(jwtManager), config.Faskes2Only())
	gin.POST("/rekam-medis-claim", handler.CreateRekamMedisandClaim)
	gin.GET("", handler.GetAllDiagnosisCodes)
}

func (h *Faskes2Handler) GetAllDiagnosisCodes(c *gin.Context) {
	infoHitter, _ := c.Get("username")
	println("Accessed by Faskes2 user:", infoHitter)

	diagnosisCodes, err := h.faskes2Service.GetAllDiagnosisCodes()
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "internal server error",
			"success": false,
			"message": "Failed to retrieve diagnosis codes",
		})
		return
	}
	c.JSON(200, gin.H{
		"data":    diagnosisCodes,
		"success": true,
		"message": "Diagnosis codes retrieved successfully",
	})
}

func (h *Faskes2Handler) CreateRekamMedisandClaim(c *gin.Context) {
	var req struct {
		RekamMedis models.RekamMedis `json:"rekam_medis"`
		Claims     models.Claims     `json:"claims"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "invalid request",
			"success": false,
			"message": "Creation Failed",
		})
		return
	}
	claimID, err := h.faskes2Service.CreateRekamMedisandClaim(req.RekamMedis, req.Claims)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "internal server error",
			"success": false,
			"message": "Creation Failed",
		})
		return
	}
	c.JSON(201, gin.H{
		"claim_id": claimID,
		"success":  true,
		"message":  "Creation Successful",
	})
}
