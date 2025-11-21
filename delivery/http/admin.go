package handlers

import (
	"net/http"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/config"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminServ models.AdminService
}

func NewAdminHandler(engine *gin.Engine, adminService models.AdminService, jwtManager *utils.JWTManager) {
	handler := &AdminHandler{adminServ: adminService}

	// Admin routes dengan auth middleware
	adminGroup := engine.Group("/admin")
	adminGroup.Use(config.AuthMiddleware(jwtManager), config.AdminOnly())

	adminGroup.GET("/claims", handler.GetAllPendingClaims)
	adminGroup.GET("/claims/:id", handler.GetClaimByID)
	adminGroup.PUT("/claims/:id/status", handler.UpdateClaimStatus)
}

func (ah *AdminHandler) GetAllPendingClaims(c *gin.Context) {
	utils.PrintContextData(c)

	claims, err := ah.adminServ.GetAllClaims(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get pending claims",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pending claims retrieved successfully",
		"data":    claims,
		"count":   len(claims),
	})
}

func (ah *AdminHandler) GetClaimByID(c *gin.Context) {
	utils.PrintContextData(c)

	claimID := c.Param("id")
	if claimID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Claim ID is required",
		})
		return
	}

	claim, err := ah.adminServ.GetClaimByID(c.Request.Context(), claimID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Claim not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Claim retrieved successfully",
		"data":    claim,
	})
}

func (ah *AdminHandler) UpdateClaimStatus(c *gin.Context) {
	utils.PrintContextData(c)

	claimID := c.Param("id")
	if claimID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Claim ID is required",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"PAID":     true,
		"REJECTED": true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid status. Allowed values: PAID and REJECTED",
		})
		return
	}

	err := ah.adminServ.UpdateClaimStatus(c.Request.Context(), claimID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update claim status",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Claim status updated successfully",
		"data": gin.H{
			"claim_id": claimID,
			"status":   req.Status,
		},
	})
}

// func (ah *AdminHandler) GetClaims(c *gin.Context) {
// 	utils.PrintContextData(c)

// 	// Optional query parameters
// 	status := c.Query("status")
// 	page := c.DefaultQuery("page", "1")
// 	limit := c.DefaultQuery("limit", "10")

// 	// Untuk sekarang, kita akan menggunakan GetAllPendingClaims
// 	// Anda bisa extend service untuk handle filtering dan pagination
// 	if status != "" && status != "PENDING" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "Only PENDING status is supported for now",
// 		})
// 		return
// 	}

// 	claims, err := ah.adminServ.GetAllPendingClaims(c.Request.Context())
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Failed to get claims",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Claims retrieved successfully",
// 		"data":    claims,
// 		"count":   len(claims),
// 		"page":    page,
// 		"limit":   limit,
// 	})
// }
