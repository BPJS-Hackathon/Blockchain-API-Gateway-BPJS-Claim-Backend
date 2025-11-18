package services

import (
	"context"
	"errors"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/domain"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo domain.AuthRepository
}

func NewAuthService(r domain.AuthRepository) *AuthService {
	return &AuthService{Repo: r}
}

func (s *AuthService) Register(ctx context.Context, username, password, role string) error {
	// check duplicate username
	_, err := s.Repo.FindByUsername(ctx, username)
	if err == nil {
		return errors.New("username already exists")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := domain.User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}

	return s.Repo.Create(&user)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.Repo.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid username or password")
	}

	// return JWT using your implemented function
	return utils.GenerateToken(user.ID, user.Role)
}
