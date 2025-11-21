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
		Preload("RekamMedis").             // Preload RekamMedis
		Preload("RekamMedis.User").        // Preload User dari RekamMedis
		Preload("RekamMedis.User.Faskes"). // Preload Faskes dari User
		Preload("RekamMedis.Diagnosis").   // Preload Diagnosis dari RekamMedis
		Order("created_at DESC").
		Find(&claims).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending claims: %v", err)
	}

	return claims, nil
}

func (r *adminRepo) GetClaimByID(ctx context.Context, clID string) (*models.Claims, error) {
	var claim models.Claims

	err := r.db.WithContext(ctx).
		Where("claim_id = ? AND status = ?", clID, "PENDING").
		Preload("RekamMedis").             // Preload RekamMedis
		Preload("RekamMedis.User").        // Preload User dari RekamMedis
		Preload("RekamMedis.User.Faskes"). // Preload Faskes dari User
		Preload("RekamMedis.Diagnosis").   // Preload Diagnosis dari RekamMedis
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
		Where("claim_id = ? AND status = ?", clID, "PENDING").
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update claim status: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("claim with ID %s not found", clID)
	}

	return nil
}

// Optional: Method untuk mendapatkan claims dengan detail lengkap
func (r *adminRepo) GetAllPendingClaimsWithFullDetails(ctx context.Context) ([]models.Claims, error) {
	var claims []models.Claims

	err := r.db.WithContext(ctx).
		Where("status = ?", "PENDING").
		Preload("RekamMedis", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Preload("User.Faskes").Preload("Diagnosis")
		}).
		Order("created_at DESC").
		Find(&claims).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending claims with full details: %v", err)
	}

	return claims, nil
}

// Optional: Method untuk mendapatkan claim dengan semua relasi
func (r *adminRepo) GetClaimWithFullDetails(ctx context.Context, clID string) (*models.Claims, error) {
	var claim models.Claims

	err := r.db.WithContext(ctx).
		Where("claim_id = ?", clID).
		Preload("RekamMedis", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Preload("User.Faskes").Preload("Diagnosis")
		}).
		First(&claim).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("claim with ID %s not found", clID)
		}
		return nil, fmt.Errorf("failed to get claim with full details: %v", err)
	}

	return &claim, nil
}
