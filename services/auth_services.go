package services

import (
	"context"
	"errors"
	"time"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/domain"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo domain.AuthRepository
	JWT  *utils.JWTManager
}

func NewAuthService(r domain.AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{
		Repo: r,
		JWT:  utils.NewJWTManager(jwtSecret, 24*time.Hour),
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.Repo.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid username or password")
	}
	token, err := s.JWT.GenerateToken(user.ID, user.Name, user.Username, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (s *AuthService) GetAccessTokenManager() *utils.JWTManager {
	return s.JWT
}

func (s *AuthService) Me(ctx context.Context, id string) (*models.User, error) {
	return s.Repo.Me(ctx, id)

}
