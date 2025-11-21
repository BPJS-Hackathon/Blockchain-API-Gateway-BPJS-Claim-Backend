package domain

import (
	"context"

	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
)

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
}

type AuthRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type AuthService interface {
	GetAccessTokenManager() *utils.JWTManager
	Login(username, password string) (string, error)
}
