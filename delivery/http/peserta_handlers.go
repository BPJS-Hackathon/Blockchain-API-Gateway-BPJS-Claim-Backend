package http

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/dto"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/services"
	"github.com/gofiber/fiber/v2"
)

type PesertaHandler struct {
	service services.PesertaService
}

func NewPesertaHandler(s services.PesertaService) *PesertaHandler {
	return &PesertaHandler{service: s}
}

func (h *PesertaHandler) CreatePeserta(c *fiber.Ctx) error {
	var req dto.PesertaBPJS
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Ambil user ID dari JWT
	userIDAny := c.Locals("user_id")
	if userIDAny == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized: missing user_id in token",
		})
	}
	userID := userIDAny.(string)
	fmt.Println(userID)

	// Convert PSTV03 (tanggal lahir)
	parsedDate, err := time.Parse("2006-01-02", req.PSTV03)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "invalid format for pstv03 (expected YYYY-MM-DD)",
		})
	}

	// Convert PSTV18 (nullable int)
	var tahunMeninggal *int
	if req.PSTV18 != nil && *req.PSTV18 != "" {
		val, convErr := strconv.Atoi(*req.PSTV18)
		if convErr != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "invalid format for pstv18 (must be int or null)",
			})
		}
		tahunMeninggal = &val
	}

	// Populate Model (struct database)
	peserta := models.PesertaBPJS{
		PSTV01:    req.PSTV01,
		PSTV02:    req.PSTV02,
		PSTV03:    parsedDate,
		PSTV04:    req.PSTV04,
		PSTV05:    req.PSTV05,
		PSTV06:    req.PSTV06,
		PSTV07:    req.PSTV07,
		PSTV08:    req.PSTV08,
		PSTV09:    req.PSTV09,
		PSTV10:    req.PSTV10,
		PSTV11:    req.PSTV11,
		PSTV12:    req.PSTV12,
		PSTV13:    req.PSTV13,
		PSTV14:    req.PSTV14,
		PSTV15:    req.PSTV15,
		PSTV16:    req.PSTV16,
		PSTV17:    req.PSTV17,
		PSTV18:    tahunMeninggal,
		CreatedBy: userID,
	}

	// Call service with the MODEL
	if err := h.service.CreatePeserta(c.Context(), &peserta); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Return saved data as DTO
	req.CreatedBy = userID
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    req,
	})
}

// GET /peserta/:nomor
func (h *PesertaHandler) GetPeserta(c *fiber.Ctx) error {
	nomor := c.Params("nomor")
	data, err := h.service.GetPesertaByNomor(c.Context(), nomor)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Peserta tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// GET /peserta/search?q=xxxx
func (h *PesertaHandler) Search(c *fiber.Ctx) error {
	keyword := c.Query("q")
	data, err := h.service.SearchPeserta(c.Context(), keyword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (h *PesertaHandler) GetAllByCreator(c *fiber.Ctx) error {
	userIDAny := c.Locals("user_id")
	if userIDAny == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized: missing user_id in token",
		})
	}

	converted := userIDAny.(string)

	datas, err := h.service.GetAllByCreator(c.Context(), converted)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    datas,
	})
}
