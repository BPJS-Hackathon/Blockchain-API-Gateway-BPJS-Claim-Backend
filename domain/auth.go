package domain

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
)

type AuthRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
}

type AuthService interface {
	GetAccessTokenManager() *utils.JWTManager
	Login(username, password string) (string, error)
}
