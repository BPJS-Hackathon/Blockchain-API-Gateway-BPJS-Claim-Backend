package services

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
)

type faskes2Service struct {
	f2rp models.Faskes2Repo
}

func NewFaskes2Service(f2rp models.Faskes2Repo) models.Faskes2Service {
	return &faskes2Service{f2rp: f2rp}
}

func (s *faskes2Service) CreateRekamMedisandClaim(ctx context.Context, rm models.RekamMedis, cl models.Claims) (string, string, error) {
	return s.f2rp.CreateRekamMedisandClaim(ctx, rm, cl)
}

func (s *faskes2Service) GetAllDiagnosisCodes(ctx context.Context) ([]models.DiagnosisCode, error) {
	return s.f2rp.GetAllDiagnosisCodes(ctx)
}

func (s *faskes2Service) GetAllPesertaDependsOnFaskesHitter(ctx context.Context, id string) (*[]models.PesertaBPJS, error) {
	return s.f2rp.GetAllPesertaDependsOnFaskesHitter(ctx, id)
}

func (s *faskes2Service) GetAllMySubmittedRMandClaim(ctx context.Context, id string) (*[]models.RekamMedis, error) {
	return s.f2rp.GetAllMySubmittedRMandClaim(ctx, id)
}
