package services

import (
	"errors"

	"emergency-rescue-locator/internal/models"
	"emergency-rescue-locator/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LocationService struct {
	locationRepo  *repositories.LocationRepository
	emergencyRepo *repositories.EmergencyRepository
}

type LocationUpdateRequest struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Accuracy  *float64 `json:"accuracy"`
}

func NewLocationService(
	locationRepo *repositories.LocationRepository,
	emergencyRepo *repositories.EmergencyRepository,
) *LocationService {
	return &LocationService{
		locationRepo:  locationRepo,
		emergencyRepo: emergencyRepo,
	}
}

func (s *LocationService) AddUpdate(emergencyID, userID uuid.UUID, req LocationUpdateRequest) (*models.LocationUpdate, error) {
	emergency, err := s.emergencyRepo.FindByID(emergencyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmergencyNotFound
		}
		return nil, err
	}

	if emergency.UserID != userID {
		return nil, ErrUnauthorized
	}

	if emergency.Status != models.StatusActive {
		return nil, ErrEmergencyNotActive
	}

	update := &models.LocationUpdate{
		EmergencyID: emergencyID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Accuracy:    req.Accuracy,
	}

	if err := s.locationRepo.Create(update); err != nil {
		return nil, err
	}

	emergency.Latitude = req.Latitude
	emergency.Longitude = req.Longitude
	_ = s.emergencyRepo.Update(emergency)

	return update, nil
}

func (s *LocationService) GetLatest(emergencyID uuid.UUID) (*models.LocationUpdate, error) {
	update, err := s.locationRepo.FindLatestByEmergencyID(emergencyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return update, nil
}

func (s *LocationService) GetHistory(emergencyID uuid.UUID) ([]models.LocationUpdate, error) {
	return s.locationRepo.FindAllByEmergencyID(emergencyID)
}
