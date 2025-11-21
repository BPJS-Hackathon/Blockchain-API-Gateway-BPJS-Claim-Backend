package services

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
)

type faskes1Service struct {
	f1rp models.Faskes1Repo
}

func NewFaskes1Service(f models.Faskes1Repo) models.Faskes1Service {
	return &faskes1Service{f1rp: f}
}

func (s *faskes1Service) CreateRekamMedis1(ctx context.Context, rm models.RekamMedis) error {
	return s.f1rp.CreateRekamMedis1(ctx, rm)
}

func (s *faskes1Service) GetAllDiagnosisCodes1(ctx context.Context) ([]models.DiagnosisCode, error) {
	return s.f1rp.GetAllDiagnosisCodes1(ctx)
}

func (s *faskes1Service) GetAllPesertaDependsOnFaskesHitter1(ctx context.Context, id string) (*[]models.PesertaBPJS, error) {
	return s.f1rp.GetAllPesertaDependsOnFaskesHitter1(ctx, id)
}
