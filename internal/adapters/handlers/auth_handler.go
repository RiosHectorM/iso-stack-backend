package handlers

import (
	"strings"
	"time"

	"github.com/RiosHectorM/iso-stack/internal/adapters/auth"
	"github.com/RiosHectorM/iso-stack/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB         *gorm.DB
	JWTAdapter *auth.JWTAdapter
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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	newOrg := domain.Organization{ID: uuid.New().String(), Name: req.OrgName}
	newUser := domain.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newOrg).Error; err != nil {
			return err
		}
		return tx.Create(&newUser).Error
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not create account"})
	}

	token, _ := h.JWTAdapter.GenerateToken(&newUser)
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

	var user domain.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	token, _ := h.JWTAdapter.GenerateToken(&user)
	return c.JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == "" {
		return c.Status(400).JSON(fiber.Map{"error": "no hay token para invalidar"})
	}

	// Guardar el token en la lista negra
	revoked := domain.RevokedToken{
		Token:     tokenString,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Tiempo de vida del JWT
	}

	if err := h.DB.Create(&revoked).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "error al cerrar sesión"})
	}

	return c.JSON(fiber.Map{"message": "sesión cerrada exitosamente"})
}
