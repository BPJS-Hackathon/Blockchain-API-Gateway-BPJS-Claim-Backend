package repositories

import (
	"context"
	"fmt"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type authRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *authRepo {
	return &authRepo{DB: db}
}

func (r *authRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var m models.User
	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&m).Error
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           m.ID,
		Name:         m.Name,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		Role:         m.Role,
	}, nil
}

func (s *authRepo) Me(ctx context.Context, id string) (*models.User, error) {
	var userData models.User

	err := s.DB.WithContext(ctx).
		Preload("Faskes"). // Preload Faskes data jika ada relasi
		Where("id = ?", id).
		First(&userData).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get user data: %v", err)
	}

	return &userData, nil
}
