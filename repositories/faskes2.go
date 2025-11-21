package repositories

import (
	"context"
	"fmt"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type faskes2Repository struct {
	db *gorm.DB
}

func NewFaskes2Repo(db *gorm.DB) models.Faskes2Repo {
	return &faskes2Repository{db: db}
}

func (r *faskes2Repository) CreateRekamMedisandClaim(ctx context.Context, rm models.RekamMedis, cl models.Claims) (string, string, error) {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(&rm).Error; err != nil {
		tx.Rollback()
		return "", "", err
	}

	cl.RekamMedisID = rm.ID
	if err := tx.Create(&cl).Error; err != nil {
		tx.Rollback()
		return "", "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", "", err
	}

	// return created claim id and nil
	return cl.ClaimID, cl.RekamMedisID, nil
}

func (r *faskes2Repository) GetAllDiagnosisCodes(ctx context.Context) ([]models.DiagnosisCode, error) {
	var diagnosisCodes []models.DiagnosisCode
	if err := r.db.WithContext(ctx).Find(&diagnosisCodes).Error; err != nil {
		return nil, err
	}
	return diagnosisCodes, nil
}

// repositories/faskes2.go - update GetAllPesertaDependsOnFaskesHitter
func (r *faskes2Repository) GetAllPesertaDependsOnFaskesHitter(ctx context.Context, userID string) (*[]models.PesertaBPJS, error) {
	var peserta []models.PesertaBPJS

	// Dapatkan informasi faskes dari user yang sedang login
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Faskes").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	if user.Faskes == nil {
		return nil, fmt.Errorf("user is not associated with any faskes")
	}

	fmt.Printf("DEBUG: User Faskes - Nama: %s, Jenis: %s\n", user.Faskes.NamaFaskes, user.Faskes.JenisFaskes)

	// Filter peserta berdasarkan nama faskes yang sama dengan faskes user
	err = r.db.WithContext(ctx).
		Where("pstv11 = ?", user.Faskes.NamaFaskes).
		Order("pstv01 ASC").
		Find(&peserta).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get peserta: %v", err)
	}

	fmt.Printf("DEBUG: Found %d peserta for faskes %s\n", len(peserta), user.Faskes.NamaFaskes)
	return &peserta, nil
}

func (r *faskes2Repository) GetAllMySubmittedRMandClaim(ctx context.Context, userID string) (*[]models.RekamMedis, error) {
	var rekamMedis []models.RekamMedis

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("User.Faskes").
		Preload("Diagnosis").
		Preload("PesertaBPJS").
		Preload("Claims").
		Order("admission_date DESC").
		Find(&rekamMedis).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get rekam medis and claims: %v", err)
	}

	fmt.Printf("DEBUG: Found %d rekam medis with claims for user %s\n", len(rekamMedis), userID)
	return &rekamMedis, nil
}
