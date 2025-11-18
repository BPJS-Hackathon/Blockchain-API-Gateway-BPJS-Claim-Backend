package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/repositories"
	"github.com/google/uuid"
)

type ClaimService interface {
	SubmitClaim(ctx context.Context, req *models.Claim, userID *string) (*models.Claim, error)
	ApproveClaim(ctx context.Context, claimID, approverID, decision, note string) error
	GetClaim(ctx context.Context, id string) (*models.Claim, error)
	ListClaims(ctx context.Context, status string, limit, offset int) ([]models.Claim, error)
}

type claimService struct {
	claimRepo repositories.ClaimRepository
	bcClient  BlockchainClient
	// approvalRepo, auditRepo can be added
}

func NewClaimService(cr repositories.ClaimRepository, bc BlockchainClient) ClaimService {
	return &claimService{claimRepo: cr, bcClient: bc}
}

func (s *claimService) SubmitClaim(ctx context.Context, req *models.Claim, userID *string) (*models.Claim, error) {
	// idempotency
	if existing, _ := s.claimRepo.GetByExternalID(ctx, req.ExternalID); existing != nil {
		return nil, errors.New("claim already exists")
	}

	// set ID if not present
	if req.ID == "" {
		req.ID = uuid.NewString()
	}
	req.Status = models.StatusPending
	if err := s.claimRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	// build tx payload
	txPayload := map[string]any{
		"type":  "CLAIM_REQUEST",
		"claim": req,
	}
	bytesPayload, _ := json.Marshal(txPayload)

	txHash, err := s.bcClient.SubmitTransaction(ctx, "CLAIM_REQUEST", bytesPayload)
	if err != nil {
		return nil, err
	}
	req.TxHash = &txHash
	req.Status = models.StatusSubmitted
	if err := s.claimRepo.Update(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *claimService) ApproveClaim(ctx context.Context, claimID, approverID, decision, note string) error {
	claim, err := s.claimRepo.GetByID(ctx, claimID)
	if err != nil {
		return err
	}
	// create approval record omitted for brevity (should write to approval repo)
	if decision == "approve" {
		claim.Status = models.StatusApproved
	} else {
		claim.Status = models.StatusRejected
	}
	if err := s.claimRepo.Update(ctx, claim); err != nil {
		return err
	}

	// submit approval tx
	txPayload := map[string]any{
		"type":        "CLAIM_APPROVAL",
		"claim_id":    claim.ID,
		"approver_id": approverID,
		"decision":    decision,
		"note":        note,
	}
	b, _ := json.Marshal(txPayload)
	txHash, err := s.bcClient.SubmitTransaction(ctx, "CLAIM_APPROVAL", b)
	if err != nil {
		return err
	}
	claim.TxHash = &txHash
	claim.Status = models.StatusSubmitted
	return s.claimRepo.Update(ctx, claim)
}

func (s *claimService) GetClaim(ctx context.Context, id string) (*models.Claim, error) {
	return s.claimRepo.GetByID(ctx, id)
}

func (s *claimService) ListClaims(ctx context.Context, status string, limit, offset int) ([]models.Claim, error) {
	return s.claimRepo.ListByStatus(ctx, status, limit, offset)
}
