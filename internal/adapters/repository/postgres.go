package repository

import (
	"fmt"
	"log"
	"github.com/RiosHectorM/iso-stack/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("No se pudo conectar a la DB")
	}

	// Auto-Migración de tablas
	err = db.AutoMigrate(&domain.Organization{}, &domain.User{}, &domain.RevokedToken{})
	if err != nil {
		log.Fatal("Error en la migración:", err)
	}

	fmt.Println("Conexión a DB y migración exitosa")
	return db
}