package controllers

import (
	"net/http"

	"emergency-rescue-locator/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmergencyController struct {
	emergencyService *services.EmergencyService
	locationService  *services.LocationService
}

func NewEmergencyController(
	emergencyService *services.EmergencyService,
	locationService *services.LocationService,
) *EmergencyController {
	return &EmergencyController{
		emergencyService: emergencyService,
		locationService:  locationService,
	}
}

func (c *EmergencyController) Create(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uuid.UUID)

	var req services.CreateEmergencyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Latitude == 0 && req.Longitude == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude are required"})
		return
	}

	emergency, err := c.emergencyService.Create(userID, req)
	if err != nil {
		if err == services.ErrActiveEmergencyExists {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create emergency"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "emergency SOS created successfully",
		"data":    emergency,
	})
}

func (c *EmergencyController) GetActive(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uuid.UUID)

	emergency, err := c.emergencyService.GetActiveByUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch active emergency"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": emergency})
}

func (c *EmergencyController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid emergency id"})
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)
	role, _ := ctx.Get("userRole")

	emergency, err := c.emergencyService.GetByID(id)
	if err != nil {
		if err == services.ErrEmergencyNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch emergency"})
		return
	}

	if role != "admin" && emergency.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	latestLocation, _ := c.locationService.GetLatest(id)

	ctx.JSON(http.StatusOK, gin.H{
		"data":            emergency,
		"latest_location": latestLocation,
	})
}

func (c *EmergencyController) Cancel(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid emergency id"})
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)

	emergency, err := c.emergencyService.Cancel(id, userID)
	if err != nil {
		switch err {
		case services.ErrEmergencyNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case services.ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case services.ErrEmergencyNotActive:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel emergency"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "emergency cancelled",
		"data":    emergency,
	})
}
