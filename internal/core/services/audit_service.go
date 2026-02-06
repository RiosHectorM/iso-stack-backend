package services

import (
	"errors"

	"github.com/RiosHectorM/iso-stack/internal/core/domain"
	"github.com/RiosHectorM/iso-stack/internal/core/ports"
	"github.com/google/uuid"
)

type AuditService struct {
	repo    ports.AuditRepository
	orgRepo ports.OrganizationRepository
}

func NewAuditService(repo ports.AuditRepository, orgRepo ports.OrganizationRepository) *AuditService {
	return &AuditService{
		repo:    repo,
		orgRepo: orgRepo,
	}
}

func (s *AuditService) CreateAudit(title, orgOwnerID, userID string) (*domain.Audit, error) {
	audit := &domain.Audit{
		Title:      title,
		OrgOwnerID: orgOwnerID,
		Status:     domain.AuditPlanificada,
	}

	if err := s.repo.CreateAudit(audit); err != nil {
		return nil, err
	}

	// Auto-assign Creator as Auditor_Lider
	assignment := &domain.AuditAssignment{
		AuditID:          audit.ID,
		UserID:           userID,
		RoleInAudit:      domain.RoleAuditorLider,
		AcceptanceStatus: domain.AcceptAceptado,
		IsActive:         true,
	}
	if err := s.repo.AssignUserToAudit(assignment); err != nil {
		return nil, err
	}

	return audit, nil
}

func (s *AuditService) AssignStaff(auditID, userID, role, orgID string) error {
	// 1. Verify User belongs to Organization
	if _, err := s.orgRepo.FindUserOrg(userID, orgID); err != nil {
		return errors.New("user does not belong to your organization")
	}

	// 2. Create Assignment
	assignment := &domain.AuditAssignment{
		AuditID:          auditID,
		UserID:           userID,
		RoleInAudit:      domain.Role(role),
		AcceptanceStatus: domain.AcceptPendiente,
		IsActive:         true,
	}

	// 3. Generate Temporary Link for External Roles
	if role == string(domain.RoleAuxiliar) || role == string(domain.RoleObservador) {
		assignment.TemporaryLink = uuid.New().String()
	}

	return s.repo.AssignUserToAudit(assignment)
}

func (s *AuditService) GetMyAudits(userID string) ([]domain.Audit, error) {
	return s.repo.GetAuditsByUserID(userID)
}

func (s *AuditService) GetPublicAudit(tempLink string) (*domain.Audit, error) {
	return s.repo.GetAuditByTempLink(tempLink)
}
