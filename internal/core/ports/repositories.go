package ports

import "github.com/RiosHectorM/iso-stack/internal/core/domain"

type AuthRepository interface {
	CreateUserWithOrg(user *domain.User, org *domain.Organization, userOrg *domain.UserOrganization) error
	FindUserByEmail(email string) (*domain.User, error)
	GetUserPrimaryOrg(userID string) (*domain.UserOrganization, error)
	RevokeToken(token string, expirationTime int64) error // expirationTime podría ser time.Time
	IsTokenRevoked(token string) (bool, error)
}

type OrganizationRepository interface {
	AddUserToOrg(userOrg *domain.UserOrganization) error
	CreateUserAndAddToOrg(user *domain.User, userOrg *domain.UserOrganization) error
	ListOrgStaff(orgID string) ([]domain.UserOrganization, error)
	FindUserOrg(userID, orgID string) (*domain.UserOrganization, error)
	UpdateUserStatus(userID, orgID string, status domain.MemberStatus) error
}

type AuditRepository interface {
	CreateAudit(audit *domain.Audit) error
	AssignUserToAudit(assignment *domain.AuditAssignment) error
	GetAuditsByUserID(userID string) ([]domain.Audit, error)
	GetAuditByTempLink(tempLink string) (*domain.Audit, error)
	GetAuditByID(auditID string) (*domain.Audit, error)
	FindAuditAssignment(auditID, userID string) (*domain.AuditAssignment, error)
}
