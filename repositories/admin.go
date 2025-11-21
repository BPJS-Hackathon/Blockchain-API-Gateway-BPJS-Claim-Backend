package repositories

import (
	"context"
	"fmt"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) models.AdminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) GetAllPendingClaims(ctx context.Context) ([]models.Claims, error) {
	var claims []models.Claims

	err := r.db.WithContext(ctx).
		Where("status = ?", "PENDING").
		Order("created_at DESC"). // Jika ada created_at, jika tidak hapus line ini
		Find(&claims).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending claims: %v", err)
	}

	return claims, nil
}

func (r *adminRepo) GetClaimByID(ctx context.Context, clID string) (*models.Claims, error) {
	var claim models.Claims

	err := r.db.WithContext(ctx).
		Where("claim_id = ?", clID).
		First(&claim).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("claim with ID %s not found", clID)
		}
		return nil, fmt.Errorf("failed to get claim: %v", err)
	}

	return &claim, nil
}

func (r *adminRepo) UpdateClaimStatus(ctx context.Context, clID string, status string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Claims{}).
		Where("claim_id = ?", clID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update claim status: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("claim with ID %s not found", clID)
	}

	return nil
}
