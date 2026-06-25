package controllers

import (
	"net/http"

	"emergency-rescue-locator/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LocationController struct {
	locationService *services.LocationService
}

func NewLocationController(locationService *services.LocationService) *LocationController {
	return &LocationController{locationService: locationService}
}

func (c *LocationController) UpdateLocation(ctx *gin.Context) {
	emergencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid emergency id"})
		return
	}

	userID := ctx.MustGet("userID").(uuid.UUID)

	var req services.LocationUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Latitude == 0 && req.Longitude == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude are required"})
		return
	}

	update, err := c.locationService.AddUpdate(emergencyID, userID, req)
	if err != nil {
		switch err {
		case services.ErrEmergencyNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case services.ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case services.ErrEmergencyNotActive:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update location"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "location updated",
		"data":    update,
	})
}

func (c *LocationController) GetLatest(ctx *gin.Context) {
	emergencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid emergency id"})
		return
	}

	update, err := c.locationService.GetLatest(emergencyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch location"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": update})
}

func (c *LocationController) GetHistory(ctx *gin.Context) {
	emergencyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid emergency id"})
		return
	}

	history, err := c.locationService.GetHistory(emergencyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch location history"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": history})
}
