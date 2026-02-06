package handlers

import (
	"strings"

	"github.com/RiosHectorM/iso-stack/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Service ports.AuthService
}

func NewAuthHandler(service ports.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		OrgName  string `json:"org_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	token, err := h.Service.Register(req.Email, req.Password, req.OrgName)
	if err != nil {
		if err.Error() == "el usuario ya existe" {
			return c.Status(409).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": "could not create account"})
	}

	return c.Status(201).JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	token, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == "" {
		return c.Status(400).JSON(fiber.Map{"error": "no hay token para invalidar"})
	}

	if err := h.Service.Logout(tokenString); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "error al cerrar sesión"})
	}

	return c.JSON(fiber.Map{"message": "sesión cerrada exitosamente"})
}
