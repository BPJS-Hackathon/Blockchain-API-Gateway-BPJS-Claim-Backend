package domain

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
)

type AuthRepository interface {
	Me(ctx context.Context, id string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
}

type AuthService interface {
	Me(ctx context.Context, id string) (*models.User, error)
	GetAccessTokenManager() *utils.JWTManager
	Login(username, password string) (string, error)
}
