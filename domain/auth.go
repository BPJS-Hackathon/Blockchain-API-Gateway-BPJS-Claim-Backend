package domain

import "context"

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
}

type AuthRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	Create(user *User) error
}

type AuthService interface {
	Register(username, password, role string) error
	Login(username, password string) (string, error)
}
