package ports

import "github.com/RiosHectorM/iso-stack/internal/core/domain"

type UserRepository interface {
	CreateUser(user *domain.User) error
	CreateOrganization(org *domain.Organization) error
	FindByEmail(email string) (*domain.User, error)
}

type AuthService interface {
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(token string) (*domain.User, error)
}