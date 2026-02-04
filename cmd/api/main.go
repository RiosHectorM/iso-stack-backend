package main

import (
	"fmt" // Importante para el print manual
	"log"

	"github.com/RiosHectorM/iso-stack/internal/adapters/auth"
	"github.com/RiosHectorM/iso-stack/internal/adapters/handlers"
	"github.com/RiosHectorM/iso-stack/internal/adapters/repository"
	"github.com/RiosHectorM/iso-stack/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()
	db := repository.NewPostgresDB(cfg.DBDSN)
	jwtAdapter := &auth.JWTAdapter{Secret: cfg.JWTSecret}

	authHandler := &handlers.AuthHandler{
		DB:         db,
		JWTAdapter: jwtAdapter,
	}

	app := fiber.New(fiber.Config{
		AppName: "ISO Stack API v1.0",
	})

	app.Use(logger.New())

	api := app.Group("/api/v1")

	// Auth
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", handlers.AuthMiddleware(cfg.JWTSecret, db), authHandler.Logout)

	// Projects (Ruta que te da 404)
	projects := api.Group("/projects")
	projects.Get("/test", handlers.AuthMiddleware(cfg.JWTSecret, db), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK", "org": c.Locals("org_id")})
	})

	// --- ESTO ES LO QUE TE VA A DECIR LA VERDAD ---
	fmt.Println("\n--- RUTAS REGISTRADAS ---")
	for _, route := range app.GetRoutes() {
		if route.Method != "USE" { // Ignorar middlewares de la lista
			fmt.Printf("%s\t%s\n", route.Method, route.Path)
		}
	}
	fmt.Println("-------------------------\n")

	log.Printf("Iniciando servidor en puerto %s...", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
