package controllers

import (
	"net/http"

	"emergency-rescue-locator/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminController struct {
	emergencyService *services.EmergencyService
	locationService  *services.LocationService
}

func NewAdminController(
	emergencyService *services.EmergencyService,
	locationService *services.LocationService,
) *AdminController {
	return &AdminController{
		emergencyService: emergencyService,
		locationService:  locationService,
	}
}

func (c *AdminController) ListActive(ctx *gin.Context) {
	emergencies, err := c.emergencyService.GetAllActive()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch emergencies"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": emergencies})
}

func (c *AdminController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	status := ctx.Query("status")

	emergencies, err := c.emergencyService.Search(query, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search emergencies"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": emergencies})
}

func (c *AdminController) GetDetails(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid emergency id"})
		return
	}

	emergency, err := c.emergencyService.GetByID(id)
	if err != nil {
		if err == services.ErrEmergencyNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch emergency"})
		return
	}

	latestLocation, _ := c.locationService.GetLatest(id)
	history, _ := c.locationService.GetHistory(id)

	ctx.JSON(http.StatusOK, gin.H{
		"data":             emergency,
		"latest_location":  latestLocation,
		"location_history": history,
	})
}

func (c *AdminController) Resolve(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid emergency id"})
		return
	}

	emergency, err := c.emergencyService.Resolve(id)
	if err != nil {
		switch err {
		case services.ErrEmergencyNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case services.ErrEmergencyNotActive:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve emergency"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "emergency marked as resolved",
		"data":    emergency,
	})
}

func (c *AdminController) Stats(ctx *gin.Context) {
	stats, err := c.emergencyService.GetStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch statistics"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": stats})
}
