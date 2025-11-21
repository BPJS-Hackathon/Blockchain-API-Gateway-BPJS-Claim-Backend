package repositories

import (
	"context"
	"fmt"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type faskes1Repository struct {
	db *gorm.DB
}

func NewFaskes1Repo(db *gorm.DB) models.Faskes1Repo {
	return &faskes1Repository{db: db}
}

// repositories/faskes1.go
func (r *faskes1Repository) GetAllMySubmittedRM1(ctx context.Context, userID string) (*[]models.RekamMedis, error) {
	var rekamMedis []models.RekamMedis

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").              // Preload User yang membuat
		Preload("User.Faskes").       // Preload Faskes dari User
		Preload("Diagnosis").         // Preload Diagnosis
		Preload("PesertaBPJS").       // Preload PesertaBPJS
		Order("admission_date DESC"). // Urutkan berdasarkan tanggal masuk terbaru
		Find(&rekamMedis).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get rekam medis: %v", err)
	}

	return &rekamMedis, nil
}

func (r *faskes1Repository) CreateRekamMedis1(ctx context.Context, rm models.RekamMedis) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(&rm).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// return created claim id and nil
	return nil
}

func (r *faskes1Repository) GetAllDiagnosisCodes1(ctx context.Context) ([]models.DiagnosisCode, error) {
	var diagnosisCodes []models.DiagnosisCode
	if err := r.db.WithContext(ctx).Find(&diagnosisCodes).Error; err != nil {
		return nil, err
	}
	return diagnosisCodes, nil
}

// repositories/faskes2.go - update GetAllPesertaDependsOnFaskesHitter
func (r *faskes1Repository) GetAllPesertaDependsOnFaskesHitter1(ctx context.Context, userID string) (*[]models.PesertaBPJS, error) {
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
