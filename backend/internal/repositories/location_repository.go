package repositories

import (
	"emergency-rescue-locator/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(update *models.LocationUpdate) error {
	return r.db.Create(update).Error
}

func (r *LocationRepository) FindLatestByEmergencyID(emergencyID uuid.UUID) (*models.LocationUpdate, error) {
	var update models.LocationUpdate
	err := r.db.Where("emergency_id = ?", emergencyID).
		Order("recorded_at DESC").
		First(&update).Error
	if err != nil {
		return nil, err
	}
	return &update, nil
}

func (r *LocationRepository) FindAllByEmergencyID(emergencyID uuid.UUID) ([]models.LocationUpdate, error) {
	var updates []models.LocationUpdate
	err := r.db.Where("emergency_id = ?", emergencyID).
		Order("recorded_at ASC").
		Find(&updates).Error
	return updates, err
}
