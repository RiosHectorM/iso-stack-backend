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

// --- OrganizationRepository Implementation ---

func (r *PostgresRepository) AddUserToOrg(userOrg *domain.UserOrganization) error {
	return r.DB.Create(userOrg).Error
}

func (r *PostgresRepository) CreateUserAndAddToOrg(user *domain.User, userOrg *domain.UserOrganization) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		userOrg.UserID = user.ID
		if err := tx.Create(userOrg).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *PostgresRepository) ListOrgStaff(orgID string) ([]domain.UserOrganization, error) {
	var members []domain.UserOrganization
	if err := r.DB.Where("organization_id = ?", orgID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *PostgresRepository) FindUserOrg(userID, orgID string) (*domain.UserOrganization, error) {
	var member domain.UserOrganization
	if err := r.DB.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *PostgresRepository) UpdateUserStatus(userID, orgID string, status domain.MemberStatus) error {
	return r.DB.Model(&domain.UserOrganization{}).
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		Update("status", status).Error
}

// --- AuditRepository Implementation ---

func (r *PostgresRepository) CreateAudit(audit *domain.Audit) error {
	return r.DB.Create(audit).Error
}

func (r *PostgresRepository) AssignUserToAudit(assignment *domain.AuditAssignment) error {
	return r.DB.Create(assignment).Error
}

func (r *PostgresRepository) GetAuditsByUserID(userID string) ([]domain.Audit, error) {
	var audits []domain.Audit
	// JOIN simple: obtener audits donde exista un assignment activo para este userID
	err := r.DB.Joins("JOIN audit_assignments ON audit_assignments.audit_id = audits.id").
		Where("audit_assignments.user_id = ? AND audit_assignments.is_active = ?", userID, true).
		Find(&audits).Error
	return audits, err
}

func (r *PostgresRepository) GetAuditByTempLink(tempLink string) (*domain.Audit, error) {
	var assignment domain.AuditAssignment
	if err := r.DB.Where("temporary_link = ? AND is_active = ?", tempLink, true).First(&assignment).Error; err != nil {
		return nil, err
	}

	// Si encontramos el assignment, devolvemos la auditoría (si no está finalizada - pending logic check)
	var audit domain.Audit
	if err := r.DB.First(&audit, "id = ?", assignment.AuditID).Error; err != nil {
		return nil, err
	}

	if audit.Status == domain.AuditFinalizada {
		return nil, fmt.Errorf("audit is finalized")
	}

	return &audit, nil
}

func (r *PostgresRepository) GetAuditByID(auditID string) (*domain.Audit, error) {
	var audit domain.Audit
	if err := r.DB.First(&audit, "id = ?", auditID).Error; err != nil {
		return nil, err
	}
	return &audit, nil
}

func (r *PostgresRepository) FindAuditAssignment(auditID, userID string) (*domain.AuditAssignment, error) {
	var assignment domain.AuditAssignment
	err := r.DB.Where("audit_id = ? AND user_id = ?", auditID, userID).First(&assignment).Error
	return &assignment, err
}
