package domain

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	Email          string         `gorm:"unique;not null" json:"email"`
	Password       string         `gorm:"not null" json:"-"` // "-" oculta el hash en los JSON
	Role           string         `gorm:"not null" json:"role"` // Admin_Org, Auditor_Lider, etc.
	OrganizationID string         `gorm:"not null" json:"org_id"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// Antes de crear en DB, generamos el UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}