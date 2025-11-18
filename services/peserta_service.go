package services

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/repositories"
)

type PesertaService interface {
	GetPesertaByNomor(ctx context.Context, nomor string) (*models.PesertaBPJS, error)
	SearchPeserta(ctx context.Context, keyword string) ([]models.PesertaBPJS, error)
	CreatePeserta(ctx context.Context, p *models.PesertaBPJS) error
	GetAllByCreator(ctx context.Context, creatorID string) ([]models.PesertaBPJS, error)
}

type pesertaService struct {
	repo repositories.PesertaRepository
}

func NewPesertaService(repo repositories.PesertaRepository) PesertaService {
	return &pesertaService{repo: repo}
}

func (s *pesertaService) GetAllByCreator(ctx context.Context, hospitalID string) ([]models.PesertaBPJS, error) {
	return s.repo.GetAllByCreator(ctx, hospitalID)
}

func (s *pesertaService) CreatePeserta(ctx context.Context, p *models.PesertaBPJS) error {
	return s.repo.Create(ctx, p)
}

func (s *pesertaService) GetPesertaByNomor(ctx context.Context, nomor string) (*models.PesertaBPJS, error) {
	return s.repo.FindByNomorPeserta(ctx, nomor)
}

func (s *pesertaService) SearchPeserta(ctx context.Context, keyword string) ([]models.PesertaBPJS, error) {
	return s.repo.Search(ctx, keyword, 20)
}
