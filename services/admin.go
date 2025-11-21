package services

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
)

type adminServ struct {
	adminRep models.AdminRepo
}

func NewAdminService(as models.AdminRepo) models.AdminService {
	return &adminServ{adminRep: as}
}

func (s *adminServ) GetAllPendingClaims(ctx context.Context) ([]models.Claims, error) {
	return s.adminRep.GetAllPendingClaims(ctx)
}

func (s *adminServ) GetClaimByID(ctx context.Context, claimID string) (*models.Claims, error) {
	return s.adminRep.GetClaimByID(ctx, claimID)
}

func (s *adminServ) UpdateClaimStatus(ctx context.Context, claimID string, status string) error {
	return s.adminRep.UpdateClaimStatus(ctx, claimID, status)
}
