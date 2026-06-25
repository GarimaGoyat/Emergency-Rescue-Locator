package services

import (
	"errors"
	"time"

	"emergency-rescue-locator/internal/models"
	"emergency-rescue-locator/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrActiveEmergencyExists = errors.New("you already have an active emergency")
	ErrEmergencyNotFound     = errors.New("emergency not found")
	ErrEmergencyNotActive    = errors.New("emergency is not active")
	ErrUnauthorized          = errors.New("unauthorized access to emergency")
)

type EmergencyService struct {
	emergencyRepo *repositories.EmergencyRepository
	locationRepo  *repositories.LocationRepository
}

type CreateEmergencyRequest struct {
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Address     string  `json:"address"`
}

type EmergencyStats struct {
	TotalEmergencies   int64 `json:"total_emergencies"`
	ActiveEmergencies  int64 `json:"active_emergencies"`
	ResolvedEmergencies int64 `json:"resolved_emergencies"`
	TodayEmergencies   int64 `json:"today_emergencies"`
}

func NewEmergencyService(
	emergencyRepo *repositories.EmergencyRepository,
	locationRepo *repositories.LocationRepository,
) *EmergencyService {
	return &EmergencyService{
		emergencyRepo: emergencyRepo,
		locationRepo:  locationRepo,
	}
}

func (s *EmergencyService) Create(userID uuid.UUID, req CreateEmergencyRequest) (*models.Emergency, error) {
	_, err := s.emergencyRepo.FindActiveByUserID(userID)
	if err == nil {
		return nil, ErrActiveEmergencyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	emergency := &models.Emergency{
		UserID:      userID,
		Status:      models.StatusActive,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Address:     req.Address,
	}

	if err := s.emergencyRepo.Create(emergency); err != nil {
		return nil, err
	}

	locationUpdate := &models.LocationUpdate{
		EmergencyID: emergency.ID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		RecordedAt:  time.Now().UTC(),
	}
	_ = s.locationRepo.Create(locationUpdate)

	return s.emergencyRepo.FindByID(emergency.ID)
}

func (s *EmergencyService) GetActiveByUser(userID uuid.UUID) (*models.Emergency, error) {
	emergency, err := s.emergencyRepo.FindActiveByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return emergency, nil
}

func (s *EmergencyService) GetByID(id uuid.UUID) (*models.Emergency, error) {
	emergency, err := s.emergencyRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmergencyNotFound
		}
		return nil, err
	}
	return emergency, nil
}

func (s *EmergencyService) GetAllActive() ([]models.Emergency, error) {
	return s.emergencyRepo.FindAllActive()
}

func (s *EmergencyService) Search(query, status string) ([]models.Emergency, error) {
	return s.emergencyRepo.Search(query, status)
}

func (s *EmergencyService) Resolve(id uuid.UUID) (*models.Emergency, error) {
	emergency, err := s.emergencyRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmergencyNotFound
		}
		return nil, err
	}

	if emergency.Status != models.StatusActive {
		return nil, ErrEmergencyNotActive
	}

	now := time.Now().UTC()
	emergency.Status = models.StatusResolved
	emergency.ResolvedAt = &now

	if err := s.emergencyRepo.Update(emergency); err != nil {
		return nil, err
	}

	return emergency, nil
}

func (s *EmergencyService) Cancel(id uuid.UUID, userID uuid.UUID) (*models.Emergency, error) {
	emergency, err := s.emergencyRepo.FindByID(id)
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

	emergency.Status = models.StatusCancelled
	if err := s.emergencyRepo.Update(emergency); err != nil {
		return nil, err
	}

	return emergency, nil
}

func (s *EmergencyService) GetStats() (*EmergencyStats, error) {
	total, err := s.emergencyRepo.CountAll()
	if err != nil {
		return nil, err
	}
	active, err := s.emergencyRepo.CountByStatus(models.StatusActive)
	if err != nil {
		return nil, err
	}
	resolved, err := s.emergencyRepo.CountByStatus(models.StatusResolved)
	if err != nil {
		return nil, err
	}
	today, err := s.emergencyRepo.CountToday()
	if err != nil {
		return nil, err
	}

	return &EmergencyStats{
		TotalEmergencies:    total,
		ActiveEmergencies:   active,
		ResolvedEmergencies: resolved,
		TodayEmergencies:    today,
	}, nil
}
