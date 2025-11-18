package repositories

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/domain"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{DB: db}
}

func (r *AuthRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m models.User
	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&m).Error
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           m.ID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		Role:         m.Role,
	}, nil
}

func (r *AuthRepo) Create(u *domain.User) error {
	m := models.User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
	}
	return r.DB.Create(&m).Error
}
