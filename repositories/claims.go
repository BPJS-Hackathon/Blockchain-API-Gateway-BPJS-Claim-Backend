// repositories/claims.go
package repositories

import (
	"context"
	"fmt"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type claimsRepository struct {
	db *gorm.DB
}

func NewClaimsRepo(db *gorm.DB) models.ClaimsRepo {
	return &claimsRepository{db: db}
}

func (r *claimsRepository) UpdateClaimStatus(ctx context.Context, claimID string, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		models.ClaimStatusSubmitted: true,
		models.ClaimStatusPending:   true,
		models.ClaimStatusPaid:      true,
		models.ClaimStatusRejected:  true,
		models.ClaimStatusFaked:     true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	result := r.db.WithContext(ctx).
		Model(&models.Claims{}).
		Where("claim_id = ?", claimID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update claim status: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("claim with ID %s not found", claimID)
	}

	return nil
}

func (r *claimsRepository) GetClaimByID(ctx context.Context, claimID string) (*models.Claims, error) {
	var claim models.Claims

	err := r.db.WithContext(ctx).
		Where("claim_id = ?", claimID).
		First(&claim).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("claim with ID %s not found", claimID)
		}
		return nil, fmt.Errorf("failed to get claim: %v", err)
	}

	return &claim, nil
}
