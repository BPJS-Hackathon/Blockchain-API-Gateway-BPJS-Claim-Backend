package handlers

import (
	"fmt"
	"net/http"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/config"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/gin-gonic/gin"
)

type Faskes1Handler struct {
	faskes1Service models.Faskes1Service
}

func NewFaskes1Handler(engine *gin.Engine, faskes1Service models.Faskes1Service, jwtManager *utils.JWTManager) {
	handler := &Faskes1Handler{faskes1Service: faskes1Service}
	// public routes
	gin := engine.Group("/faskes1")
	gin.Use(config.AuthMiddleware(jwtManager), config.Faskes1Only())
	gin.POST("/rekam-medis", handler.CreateRekamMedis1)
	gin.GET("", handler.GetAllDiagnosisCodes1)
	gin.GET("/peserta", handler.GetAllDependedPeserta1)
}

func (h *Faskes1Handler) GetAllDiagnosisCodes1(c *gin.Context) {
	infoHitter, _ := c.Get("username")
	println("Accessed by Faskes2 user:", infoHitter)

	diagnosisCodes, err := h.faskes1Service.GetAllDiagnosisCodes1(c.Request.Context())
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

func (h *Faskes1Handler) GetAllDependedPeserta1(c *gin.Context) {
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

	data, err := h.faskes1Service.GetAllPesertaDependsOnFaskesHitter1(c.Request.Context(), userHitterID.(string))
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

func (h *Faskes1Handler) CreateRekamMedis1(c *gin.Context) {
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

	err := h.faskes1Service.CreateRekamMedis1(c.Request.Context(), req.RekamMedis)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "Creation Failed",
		})
		return
	}
	c.JSON(201, gin.H{
		"success": true,
		"message": "Creation Successful",
	})
}
