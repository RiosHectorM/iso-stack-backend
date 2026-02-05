package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- TIPOS Y ESTADOS (ENUMS) ---

type Role string

const (
	RoleConsultora     Role = "Consultora"
	RoleAuditorLider   Role = "Auditor_Lider"
	RoleAuditorInterno Role = "Auditor_Interno"
	RoleAuxiliar       Role = "Auxiliar"
	RoleObservador     Role = "Observador"
)

type MemberStatus string

const (
	MemberActivo   MemberStatus = "Activo"
	MemberInactivo MemberStatus = "Inactivo"
	MemberInvitado MemberStatus = "Invitado"
)

type AcceptanceStatus string

const (
	AcceptPendiente AcceptanceStatus = "Pendiente"
	AcceptAceptado  AcceptanceStatus = "Aceptado"
	AcceptRechazado AcceptanceStatus = "Rechazado"
)

type AuditStatus string

const (
	AuditPlanificada AuditStatus = "Planificada"
	AuditEnCurso     AuditStatus = "En_Curso"
	AuditFinalizada  AuditStatus = "Finalizada"
	AuditPausada     AuditStatus = "Pausada"
)

// --- MODELOS DE BASE DE DATOS ---

type Organization struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserOrganization es el STAFF de la empresa
type UserOrganization struct {
	UserID         string       `gorm:"primaryKey" json:"user_id"`
	OrganizationID string       `gorm:"primaryKey" json:"org_id"`
	RoleDefault    Role         `gorm:"not null" json:"role_default"`
	Status         MemberStatus `gorm:"default:'Invitado'" json:"status"`
	JoinedAt       time.Time    `json:"joined_at"`
}

type Audit struct {
	ID         string      `gorm:"primaryKey" json:"id"`
	Title      string      `gorm:"not null" json:"title"`
	OrgOwnerID string      `gorm:"not null" json:"org_owner_id"` // La empresa que la creó
	Status     AuditStatus `gorm:"default:'Planificada'" json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type AuditAssignment struct {
	AuditID          string           `gorm:"primaryKey" json:"audit_id"`
	UserID           string           `gorm:"primaryKey" json:"user_id"`
	RoleInAudit      Role             `gorm:"not null" json:"role_in_audit"`
	AcceptanceStatus AcceptanceStatus `gorm:"default:'Pendiente'" json:"acceptance_status"`
	IsActive         bool             `gorm:"default:true" json:"is_active"`
	TemporaryLink    string           `gorm:"index" json:"temporary_link,omitempty"`
}

type RevokedToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"index;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

// --- HOOKS (Generación de UUIDs) ---

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New().String()
	return
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

func (a *Audit) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New().String()
	return
}
