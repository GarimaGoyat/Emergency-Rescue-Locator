package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmergencyStatus string

const (
	StatusActive   EmergencyStatus = "active"
	StatusResolved EmergencyStatus = "resolved"
	StatusCancelled EmergencyStatus = "cancelled"
)

type Emergency struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	User        User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status      EmergencyStatus `gorm:"type:varchar(20);default:'active';index" json:"status"`
	Description string          `gorm:"type:text" json:"description"`
	Latitude    float64         `gorm:"not null" json:"latitude"`
	Longitude   float64         `gorm:"not null" json:"longitude"`
	Address     string          `gorm:"type:text" json:"address"`
	ResolvedAt  *time.Time      `json:"resolved_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (e *Emergency) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
