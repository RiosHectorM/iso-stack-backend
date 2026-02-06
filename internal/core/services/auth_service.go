package services

import (
	"errors"
	"time"

	"github.com/RiosHectorM/iso-stack/internal/adapters/auth"
	"github.com/RiosHectorM/iso-stack/internal/core/domain"
	"github.com/RiosHectorM/iso-stack/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       ports.AuthRepository
	jwtAdapter *auth.JWTAdapter
}

func NewAuthService(repo ports.AuthRepository, jwtAdapter *auth.JWTAdapter) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtAdapter: jwtAdapter,
	}
}

func (s *AuthService) Register(email, password, orgName string) (string, error) {
	// Verificar si el usuario ya existe
	if _, err := s.repo.FindUserByEmail(email); err == nil {
		return "", errors.New("el usuario ya existe")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	newOrg := &domain.Organization{Name: orgName}
	newUser := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	userOrg := &domain.UserOrganization{
		RoleDefault: domain.RoleConsultora,
		Status:      domain.MemberActivo,
	}

	// Transacción en repositorio
	if err := s.repo.CreateUserWithOrg(newUser, newOrg, userOrg); err != nil {
		return "", err
	}

	// Generar Token con contexto (OrgID creada, Rol Default)
	return s.jwtAdapter.GenerateToken(newUser.ID, newOrg.ID, string(domain.RoleConsultora))
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("credenciales inválidas")
	}

	// Obtener la organización principal del usuario
	userOrg, err := s.repo.GetUserPrimaryOrg(user.ID)
	if err != nil {
		// En un caso real podríamos devolver un token "sin org" o error.
		// Asumimos error para forzar al usuario a tener organización.
		return "", errors.New("error recuperando datos de organización del usuario")
	}

	return s.jwtAdapter.GenerateToken(user.ID, userOrg.OrganizationID, string(userOrg.RoleDefault))
}

func (s *AuthService) Logout(token string) error {
	// 24 hours expiration default
	expiration := time.Now().Add(24 * time.Hour).Unix()
	return s.repo.RevokeToken(token, expiration)
}
