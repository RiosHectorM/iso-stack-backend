package main

import (
	"fmt"
	"log"

	"github.com/RiosHectorM/iso-stack/internal/adapters/auth"
	"github.com/RiosHectorM/iso-stack/internal/adapters/handlers"
	"github.com/RiosHectorM/iso-stack/internal/adapters/repository"
	"github.com/RiosHectorM/iso-stack/internal/config"
	"github.com/RiosHectorM/iso-stack/internal/core/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()

	// 1. Adapters (Repository & Auth)
	// NewPostgresDB returns *PostgresRepository which implements ports.AuthRepository
	repo := repository.NewPostgresDB(cfg.DBDSN)
	jwtAdapter := &auth.JWTAdapter{Secret: cfg.JWTSecret}

	// 2. Application Core (Services)
	authService := services.NewAuthService(repo, jwtAdapter)

	// 3. Adapters (Handlers)
	authHandler := handlers.NewAuthHandler(authService)

	// 4. Fiber App Setup
	app := fiber.New(fiber.Config{
		AppName: "ISO Stack API v1.0",
	})

	app.Use(logger.New())

	api := app.Group("/api/v1")

	// Auth Routes
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", handlers.AuthMiddleware(cfg.JWTSecret, repo), authHandler.Logout)

	// Protected Routes (Example)
	projects := api.Group("/projects")
	projects.Get("/test", handlers.AuthMiddleware(cfg.JWTSecret, repo), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"org":     c.Locals("org_id"),
			"role":    c.Locals("role"),
			"user_id": c.Locals("user_id"),
		})
	})

	// Debug Routes Info
	fmt.Println("\n--- RUTAS REGISTRADAS ---")
	for _, route := range app.GetRoutes() {
		if route.Method != "USE" {
			fmt.Printf("%s\t%s\n", route.Method, route.Path)
		}
	}

	log.Printf("Iniciando servidor en puerto %s...", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
