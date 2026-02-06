package handlers

import (
	"github.com/RiosHectorM/iso-stack/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type AuditHandler struct {
	service ports.AuditService
}

func NewAuditHandler(service ports.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) CreateAudit(c *fiber.Ctx) error {
	var req struct {
		Title string `json:"title"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	orgID := c.Locals("org_id").(string)
	userID := c.Locals("user_id").(string)

	audit, err := h.service.CreateAudit(req.Title, orgID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(audit)
}

func (h *AuditHandler) AssignStaff(c *fiber.Ctx) error {
	var req struct {
		UserID      string `json:"user_id"`
		RoleInAudit string `json:"role_in_audit"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	auditID := c.Params("audit_id")
	orgID := c.Locals("org_id").(string)

	if err := h.service.AssignStaff(auditID, req.UserID, req.RoleInAudit, orgID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "staff assigned"})
}

func (h *AuditHandler) GetMyAudits(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	audits, err := h.service.GetMyAudits(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(audits)
}

func (h *AuditHandler) GetPublicAudit(c *fiber.Ctx) error {
	tempLink := c.Params("temp_link")
	audit, err := h.service.GetPublicAudit(tempLink)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "invalid or expired link"})
	}
	return c.JSON(audit)
}
