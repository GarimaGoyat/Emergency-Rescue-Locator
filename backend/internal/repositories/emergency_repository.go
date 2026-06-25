package repositories

import (
	"emergency-rescue-locator/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmergencyRepository struct {
	db *gorm.DB
}

func NewEmergencyRepository(db *gorm.DB) *EmergencyRepository {
	return &EmergencyRepository{db: db}
}

func (r *EmergencyRepository) Create(emergency *models.Emergency) error {
	return r.db.Create(emergency).Error
}

func (r *EmergencyRepository) FindByID(id uuid.UUID) (*models.Emergency, error) {
	var emergency models.Emergency
	err := r.db.Preload("User").First(&emergency, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &emergency, nil
}

func (r *EmergencyRepository) FindActiveByUserID(userID uuid.UUID) (*models.Emergency, error) {
	var emergency models.Emergency
	err := r.db.Where("user_id = ? AND status = ?", userID, models.StatusActive).
		Order("created_at DESC").
		First(&emergency).Error
	if err != nil {
		return nil, err
	}
	return &emergency, nil
}

func (r *EmergencyRepository) FindAllActive() ([]models.Emergency, error) {
	var emergencies []models.Emergency
	err := r.db.Preload("User").
		Where("status = ?", models.StatusActive).
		Order("created_at DESC").
		Find(&emergencies).Error
	return emergencies, err
}

func (r *EmergencyRepository) Search(query string, status string) ([]models.Emergency, error) {
	var emergencies []models.Emergency
	db := r.db.Preload("User").Order("created_at DESC")

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if query != "" {
		search := "%" + query + "%"
		db = db.Joins("JOIN users ON users.id = emergencies.user_id").
			Where(
				"users.name ILIKE ? OR users.email ILIKE ? OR emergencies.description ILIKE ? OR emergencies.address ILIKE ?",
				search, search, search, search,
			)
	}

	err := db.Find(&emergencies).Error
	return emergencies, err
}

func (r *EmergencyRepository) Update(emergency *models.Emergency) error {
	return r.db.Save(emergency).Error
}

func (r *EmergencyRepository) CountByStatus(status models.EmergencyStatus) (int64, error) {
	var count int64
	err := r.db.Model(&models.Emergency{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *EmergencyRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Emergency{}).Count(&count).Error
	return count, err
}

func (r *EmergencyRepository) CountToday() (int64, error) {
	var count int64
	err := r.db.Model(&models.Emergency{}).
		Where("created_at >= CURRENT_DATE").
		Count(&count).Error
	return count, err
}
