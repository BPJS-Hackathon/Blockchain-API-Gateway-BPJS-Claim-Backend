package repositories

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"gorm.io/gorm"
)

type PesertaRepository interface {
	FindAllByHospital(ctx context.Context, hospitalID string) ([]models.PesertaBPJS, error)
	FindByNomorPeserta(ctx context.Context, nomor string) (*models.PesertaBPJS, error)
	Search(ctx context.Context, keyword string, limit int) ([]models.PesertaBPJS, error)
	GetAllByCreator(ctx context.Context, creatorID string) ([]models.PesertaBPJS, error)
	Create(ctx context.Context, p *models.PesertaBPJS) error
}

type pesertaRepo struct {
	db *gorm.DB
}

func NewPesertaRepository(db *gorm.DB) PesertaRepository {
	return &pesertaRepo{db: db}
}

func (r *pesertaRepo) Create(ctx context.Context, p *models.PesertaBPJS) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pesertaRepo) GetAllByCreator(ctx context.Context, creatorID string) ([]models.PesertaBPJS, error) {
	var list []models.PesertaBPJS
	err := r.db.WithContext(ctx).
		Where("created_by = ?", creatorID).
		Find(&list).Error
	return list, err
}

func (r *pesertaRepo) FindByNomorPeserta(ctx context.Context, nomor string) (*models.PesertaBPJS, error) {
	var peserta models.PesertaBPJS
	err := r.db.WithContext(ctx).Where("pstv01 = ?", nomor).First(&peserta).Error
	if err != nil {
		return nil, err
	}
	return &peserta, nil
}

func (r *pesertaRepo) Search(ctx context.Context, keyword string, limit int) ([]models.PesertaBPJS, error) {
	var results []models.PesertaBPJS
	err := r.db.WithContext(ctx).Where("pstv01 ILIKE ? OR pstv02 ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Limit(limit).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *pesertaRepo) FindAllByHospital(ctx context.Context, hospitalID string) ([]models.PesertaBPJS, error) {
	var list []models.PesertaBPJS
	err := r.db.WithContext(ctx).
		Where("created_by = ?", hospitalID).
		Find(&list).Error
	return list, err
}
