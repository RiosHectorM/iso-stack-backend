package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN     string
	JWTSecret string
	Port      string
}

func LoadConfig() *Config {
	// Intentar cargar .env pero no fallar si no existe (útil para Docker/Prod)
	_ = godotenv.Load()

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("FATAL: JWT_SECRET no definido en el entorno")
	}

	return &Config{
		DBDSN:     dsn,
		JWTSecret: secret,
		Port:      os.Getenv("PORT"),
	}
}