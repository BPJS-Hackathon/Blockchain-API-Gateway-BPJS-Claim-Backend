package repositories

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type faskes2Repository struct {
	db *gorm.DB
}

func NewFaskes2Repo(db *gorm.DB) models.Faskes2Repo {
	return &faskes2Repository{db: db}
}

func (r *faskes2Repository) CreateRekamMedisandClaim(ctx context.Context, rm models.RekamMedis, cl models.Claims) (string, error) {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(&rm).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	cl.RekamMedisID = rm.ID
	if err := tx.Create(&cl).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	// return created claim id and nil
	return cl.ClaimID, nil
}

func (r *faskes2Repository) GetAllDiagnosisCodes(ctx context.Context) ([]models.DiagnosisCode, error) {
	var diagnosisCodes []models.DiagnosisCode
	if err := r.db.WithContext(ctx).Find(&diagnosisCodes).Error; err != nil {
		return nil, err
	}
	return diagnosisCodes, nil
}
