package http

import (
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClaimHandler struct {
	svc services.ClaimService
}

func NewClaimHandler(s services.ClaimService) *ClaimHandler {
	return &ClaimHandler{svc: s}
}

type SubmitClaimReq struct {
	ExternalID string                 `json:"external_id"`
	FacilityID string                 `json:"facility_id"`
	PatientID  string                 `json:"patient_id"`
	Amount     int64                  `json:"amount"`
	Payload    map[string]interface{} `json:"payload"`
}

func (h *ClaimHandler) SubmitClaim(c *fiber.Ctx) error {
	var req SubmitClaimReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	claim := &models.Claim{
		ID:         uuid.NewString(),
		FacilityID: req.FacilityID,
		PatientID:  req.PatientID,
		Amount:     req.Amount,
	}
	res, err := h.svc.SubmitClaim(c.Context(), claim, nil)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

type ApproveReq struct {
	ClaimID  string `json:"claim_id"`
	Decision string `json:"decision"`
	Note     string `json:"note"`
}

func (h *ClaimHandler) ApproveClaim(c *fiber.Ctx) error {
	var req ApproveReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := h.svc.ApproveClaim(c.Context(), req.ClaimID, userID.(string), req.Decision, req.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true})
}

func (h *ClaimHandler) GetClaim(c *fiber.Ctx) error {
	id := c.Params("id")
	claim, err := h.svc.GetClaim(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(claim)
}

func (h *ClaimHandler) ListClaims(c *fiber.Ctx) error {
	status := c.Query("status", "")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	list, err := h.svc.ListClaims(c.Context(), status, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}
