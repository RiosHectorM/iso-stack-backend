package handlers

import (
	"github.com/RiosHectorM/iso-stack/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type OrganizationHandler struct {
	service ports.OrganizationService
}

func NewOrganizationHandler(service ports.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) InviteStaff(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	orgID := c.Locals("org_id").(string)

	if err := h.service.InviteStaff(req.Email, req.Role, orgID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "invitation sent"})
}

func (h *OrganizationHandler) ListStaff(c *fiber.Ctx) error {
	orgID := c.Locals("org_id").(string)

	staff, err := h.service.ListStaff(orgID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(staff)
}

func (h *OrganizationHandler) UpdateStaffStatus(c *fiber.Ctx) error {
	var req struct {
		UserID string `json:"user_id"`
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	orgID := c.Locals("org_id").(string)

	if err := h.service.UpdateStaffStatus(req.UserID, orgID, req.Status); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "status updated"})
}
