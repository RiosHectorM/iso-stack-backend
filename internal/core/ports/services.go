package ports

import "github.com/RiosHectorM/iso-stack/internal/core/domain"

type AuthService interface {
	Register(email, password, orgName string) (string, error)
	Login(email, password string) (string, error)
	Logout(token string) error
}

type OrganizationService interface {
	InviteStaff(email, role, orgID string) error
	ListStaff(orgID string) ([]map[string]interface{}, error)
	UpdateStaffStatus(userID, orgID, status string) error
}

type AuditService interface {
	CreateAudit(title, orgOwnerID, userID string) (*domain.Audit, error)
	AssignStaff(auditID, userID, role, orgID string) error
	GetMyAudits(userID string) ([]domain.Audit, error)
	GetPublicAudit(tempLink string) (*domain.Audit, error)
}
