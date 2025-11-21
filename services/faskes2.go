package services

import "github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"

type faskes2Service struct {
	f2rp models.Faskes2Repo
}

func NewFaskes2Service(f2rp models.Faskes2Repo) models.Faskes2Service {
	return &faskes2Service{f2rp: f2rp}
}

func (s *faskes2Service) CreateRekamMedisandClaim(rm models.RekamMedis, cl models.Claims) (string, error) {
	return s.f2rp.CreateRekamMedisandClaim(rm, cl)
}

func (s *faskes2Service) GetAllDiagnosisCodes() ([]models.DiagnosisCode, error) {
	return s.f2rp.GetAllDiagnosisCodes()
}
