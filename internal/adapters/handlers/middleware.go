package handlers

import (
	"strings"

	"github.com/RiosHectorM/iso-stack/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string, repo ports.AuthRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Obtener el Header Authorization: Bearer <token>
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "falta token de sesión"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Parsear y Validar el Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "token inválido o expirado"})
		}

		// 3. Extraer Claims y Guardar en el Contexto Local
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "error al procesar claims"})
		}

		// 4. Verificar Revocación en BD
		revoked, err := repo.IsTokenRevoked(tokenString)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "error interno verificando sesión"})
		}
		if revoked {
			return c.Status(401).JSON(fiber.Map{"error": "token revocado, por favor inicie sesión nuevamente"})
		}

		// Guardamos los datos para que los handlers de negocio (Auditorías) los usen
		c.Locals("user_id", claims["user_id"])
		c.Locals("org_id", claims["org_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
