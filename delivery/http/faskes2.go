package handlers

import (
	"fmt"
	"net/http"

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
	gin.GET("/peserta", handler.GetAllDependedPeserta)
	gin.GET("/rekam-medis-claim", handler.GetAllMySubmittedRMandClaim)
}

// handlers/faskes2.go
func (h *Faskes2Handler) GetAllMySubmittedRMandClaim(c *gin.Context) {
	utils.PrintContextData(c)

	userHitterID, isVool := c.Get("userID")
	if !isVool {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"success": false,
			"message": "User ID not found in context",
		})
		return
	}

	// Panggil service untuk mendapatkan rekam medis dan claims
	rekamMedis, err := h.faskes2Service.GetAllMySubmittedRMandClaim(c.Request.Context(), userHitterID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "Failed to retrieve rekam medis and claims",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    rekamMedis,
		"success": true,
		"message": "Rekam medis and claims retrieved successfully",
		"count":   len(*rekamMedis),
	})
}

func (h *Faskes2Handler) GetAllDiagnosisCodes(c *gin.Context) {
	infoHitter, _ := c.Get("username")
	println("Accessed by Faskes2 user:", infoHitter)

	diagnosisCodes, err := h.faskes2Service.GetAllDiagnosisCodes(c.Request.Context())
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

func (h *Faskes2Handler) GetAllDependedPeserta(c *gin.Context) {
	infoHitter, _ := c.Get("username")
	userHitterID, isVool := c.Get("userID")
	if !isVool {
		c.JSON(500, gin.H{
			"error":   "Unathorized",
			"success": false,
			"message": "Fetch Failed",
		})
		return
	}
	println("Accessed by Faskes2 user:", infoHitter)

	data, err := h.faskes2Service.GetAllPesertaDependsOnFaskesHitter(c.Request.Context(), userHitterID.(string))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "Failed to retrieve peserta ",
		})
		return
	}
	c.JSON(200, gin.H{
		"data":    data,
		"success": true,
		"message": "Peserta retrieved successfully",
	})
}

func (h *Faskes2Handler) CreateRekamMedisandClaim(c *gin.Context) {
	utils.PrintContextData(c)

	userHitterID, isVool := c.Get("userID")
	if !isVool {
		c.JSON(500, gin.H{
			"error":   "Unauthorized",
			"success": false,
			"message": "Creation Failed",
		})
		return
	}

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

	// Validate Jenis Rawat
	validJenisRawat := map[string]bool{
		models.RawatInap:  true,
		models.RawatJalan: true,
	}

	if !validJenisRawat[req.RekamMedis.JenisRawat] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Creation Failed",
			"success": false,
			"error":   fmt.Sprintf("Invalid jenis rawat. Allowed values: %s, %s", models.RawatInap, models.RawatJalan),
		})
		return
	}

	// Validate Outcome (jika ada)
	if req.RekamMedis.Outcome != nil {
		validOutcome := map[string]bool{
			models.OutcomeMeninggal:   true,
			models.OutcomePulangPaksa: true,
			models.OutcomeRujuk:       true,
			models.OutcomeSembuh:      true,
		}

		if !validOutcome[*req.RekamMedis.Outcome] {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Creation Failed",
				"success": false,
				"error": fmt.Sprintf("Invalid outcome. Allowed values: %s, %s, %s, %s",
					models.OutcomeMeninggal, models.OutcomePulangPaksa, models.OutcomeRujuk, models.OutcomeSembuh),
			})
			return
		}
	}

	// Set default values
	req.RekamMedis.UserID = userHitterID.(string)

	// Set claim status jika kosong
	if req.Claims.Status == "" {
		req.Claims.Status = models.ClaimStatusPending // atau models.ClaimStatusSubmitted
	}

	claimID, err := h.faskes2Service.CreateRekamMedisandClaim(c.Request.Context(), req.RekamMedis, req.Claims)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   err.Error(),
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
