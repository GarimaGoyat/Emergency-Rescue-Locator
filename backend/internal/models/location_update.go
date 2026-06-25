package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocationUpdate struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	EmergencyID uuid.UUID      `gorm:"type:uuid;not null;index" json:"emergency_id"`
	Emergency   Emergency      `gorm:"foreignKey:EmergencyID" json:"-"`
	Latitude    float64        `gorm:"not null" json:"latitude"`
	Longitude   float64        `gorm:"not null" json:"longitude"`
	Accuracy    *float64       `json:"accuracy,omitempty"`
	RecordedAt  time.Time      `gorm:"not null;index" json:"recorded_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (l *LocationUpdate) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	if l.RecordedAt.IsZero() {
		l.RecordedAt = time.Now().UTC()
	}
	return nil
}
