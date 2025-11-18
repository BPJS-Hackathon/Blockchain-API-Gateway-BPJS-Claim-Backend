package repositories

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type ClaimRepository interface {
	Create(ctx context.Context, c *models.Claim) error
	GetByID(ctx context.Context, id string) (*models.Claim, error)
	GetByExternalID(ctx context.Context, externalID string) (*models.Claim, error)
	Update(ctx context.Context, c *models.Claim) error
	FindByTxHash(ctx context.Context, txHash string) (*models.Claim, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]models.Claim, error)
}

type claimRepo struct {
	db *gorm.DB
}

func NewClaimRepository(db *gorm.DB) ClaimRepository { return &claimRepo{db} }

func (r *claimRepo) Create(ctx context.Context, c *models.Claim) error {
	return r.db.WithContext(ctx).Create(c).Error
}
func (r *claimRepo) GetByID(ctx context.Context, id string) (*models.Claim, error) {
	var c models.Claim
	if err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *claimRepo) GetByExternalID(ctx context.Context, externalID string) (*models.Claim, error) {
	var c models.Claim
	if err := r.db.WithContext(ctx).First(&c, "external_id = ?", externalID).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *claimRepo) Update(ctx context.Context, c *models.Claim) error {
	return r.db.WithContext(ctx).Save(c).Error
}
func (r *claimRepo) FindByTxHash(ctx context.Context, txHash string) (*models.Claim, error) {
	var c models.Claim
	if err := r.db.WithContext(ctx).First(&c, "tx_hash = ?", txHash).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *claimRepo) ListByStatus(ctx context.Context, status string, limit, offset int) ([]models.Claim, error) {
	var out []models.Claim
	q := r.db.WithContext(ctx).Model(&models.Claim{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
