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
	// NewPostgresDB returns *PostgresRepository which implements ports.AuthRepository, OrgRepo, AuditRepo
	repo := repository.NewPostgresDB(cfg.DBDSN)
	jwtAdapter := &auth.JWTAdapter{Secret: cfg.JWTSecret}

	// 2. Application Core (Services)
	authService := services.NewAuthService(repo, jwtAdapter)
	orgService := services.NewOrganizationService(repo, repo) // Repo implements both interfaces
	auditService := services.NewAuditService(repo, repo)

	// 3. Adapters (Handlers)
	authHandler := handlers.NewAuthHandler(authService)
	orgHandler := handlers.NewOrganizationHandler(orgService)
	auditHandler := handlers.NewAuditHandler(auditService)

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

	// Organization Staff Routes
	orgGroup := api.Group("/organization")
	orgGroup.Use(handlers.AuthMiddleware(cfg.JWTSecret, repo))
	orgGroup.Post("/staff/invite", orgHandler.InviteStaff)
	orgGroup.Get("/staff", orgHandler.ListStaff)
	orgGroup.Patch("/staff/status", orgHandler.UpdateStaffStatus)

	// Audit Routes
	auditGroup := api.Group("/audits")
	auditGroup.Use(handlers.AuthMiddleware(cfg.JWTSecret, repo))
	auditGroup.Post("/", auditHandler.CreateAudit)
	auditGroup.Post("/:audit_id/assign", auditHandler.AssignStaff)

	// Project Routes
	projectGroup := api.Group("/projects")
	projectGroup.Use(handlers.AuthMiddleware(cfg.JWTSecret, repo))
	projectGroup.Get("/my-audits", auditHandler.GetMyAudits)

	// Public Access
	api.Get("/public/access/:temp_link", auditHandler.GetPublicAudit)

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
