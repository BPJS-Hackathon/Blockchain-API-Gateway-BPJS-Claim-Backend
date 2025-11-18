package repositories

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, u *models.User) error
}

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepo{db} }

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).First(&u, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}
