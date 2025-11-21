package handlers

import (
	"fmt"
	"net/http"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/gin-gonic/gin"
)

type BlockchainHandler struct {
	claimsRepo models.ClaimsRepo
}

func NewBlockchainHandler(engine *gin.Engine, claimsRepo models.ClaimsRepo) {
	handler := &BlockchainHandler{claimsRepo: claimsRepo}

	// Public routes untuk blockchain node - CORS khusus
	blockchainGroup := engine.Group("/blockchain")

	// Custom CORS middleware untuk blockchain
	blockchainGroup.Use(blockchainCORS())

	blockchainGroup.PUT("/claims/:id/status", handler.UpdateClaimStatus)
	blockchainGroup.GET("/claims/:id", handler.GetClaimStatus)
}

// Custom CORS middleware untuk blockchain node
func blockchainCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow specific blockchain node origins
		allowedOrigins := []string{
			"https://blockchain-bpjs.com",
			"http://localhost:3000", // untuk development
			"http://127.0.0.1:3000",
		}

		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Blockchain-Node-ID")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// UpdateClaimStatus - Endpoint untuk blockchain node update status claim
func (h *BlockchainHandler) UpdateClaimStatus(c *gin.Context) {
	// Log request info untuk audit
	fmt.Printf("=== BLOCKCHAIN NODE REQUEST ===\n")
	fmt.Printf("Method: %s\n", c.Request.Method)
	fmt.Printf("Path: %s\n", c.Request.URL.Path)
	fmt.Printf("Origin: %s\n", c.Request.Header.Get("Origin"))
	fmt.Printf("Node-ID: %s\n", c.Request.Header.Get("X-Blockchain-Node-ID"))
	fmt.Printf("==============================\n")

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
		TxHash string `json:"tx_hash"` // Optional: hash transaksi blockchain
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
		models.ClaimStatusPending:  true,
		models.ClaimStatusRejected: true,
		models.ClaimStatusFaked:    true,
		models.ClaimStatusPaid:     true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid status for blockchain update",
			"error": fmt.Sprintf("Allowed statuses: %s, %s, %s, %s",
				models.ClaimStatusPending, models.ClaimStatusRejected,
				models.ClaimStatusFaked, models.ClaimStatusPaid),
		})
		return
	}

	// Update claim status
	err := h.claimsRepo.UpdateClaimStatus(c.Request.Context(), claimID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update claim status",
			"error":   err.Error(),
		})
		return
	}

	// Log successful update
	fmt.Printf("✅ Claim %s updated to status: %s\n", claimID, req.Status)
	if req.TxHash != "" {
		fmt.Printf("📝 Transaction Hash: %s\n", req.TxHash)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Claim status updated successfully",
		"claim_id": claimID,
		"status":   req.Status,
		"tx_hash":  req.TxHash,
	})
}

// GetClaimStatus - Endpoint untuk blockchain node get claim status
func (h *BlockchainHandler) GetClaimStatus(c *gin.Context) {
	claimID := c.Param("id")
	if claimID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Claim ID is required",
		})
		return
	}

	claim, err := h.claimsRepo.GetClaimByID(c.Request.Context(), claimID)
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
