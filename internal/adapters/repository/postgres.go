package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/RiosHectorM/iso-stack/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresDB(dsn string) *PostgresRepository {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("No se pudo conectar a la DB")
	}

	// Auto-Migración de tablas
	err = db.AutoMigrate(
		&domain.Organization{},
		&domain.User{},
		&domain.UserOrganization{},
		&domain.Audit{},
		&domain.AuditAssignment{},
		&domain.RevokedToken{},
	)
	if err != nil {
		log.Fatal("Error en la migración:", err)
	}

	fmt.Println("Conexión a DB y migración exitosa")
	return &PostgresRepository{DB: db}
}

// --- AuthRepository Implementation ---

func (r *PostgresRepository) CreateUserWithOrg(user *domain.User, org *domain.Organization, userOrg *domain.UserOrganization) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(org).Error; err != nil {
			return err
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		// Asignar IDs generados al vínculo
		userOrg.UserID = user.ID
		userOrg.OrganizationID = org.ID
		if err := tx.Create(userOrg).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *PostgresRepository) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) GetUserPrimaryOrg(userID string) (*domain.UserOrganization, error) {
	var userOrg domain.UserOrganization
	// Finds the first organization (simplification for Primary)
	if err := r.DB.Where("user_id = ?", userID).First(&userOrg).Error; err != nil {
		return nil, err
	}
	return &userOrg, nil
}

func (r *PostgresRepository) RevokeToken(token string, expirationTime int64) error {
	revoked := domain.RevokedToken{
		Token:     token,
		ExpiresAt: time.Unix(expirationTime, 0),
	}
	return r.DB.Create(&revoked).Error
}

func (r *PostgresRepository) IsTokenRevoked(token string) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.RevokedToken{}).Where("token = ?", token).Count(&count).Error
	return count > 0, err
}
