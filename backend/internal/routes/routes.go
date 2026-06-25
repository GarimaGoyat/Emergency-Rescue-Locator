package routes

import (
	"net/http"

	"emergency-rescue-locator/internal/config"
	"emergency-rescue-locator/internal/controllers"
	"emergency-rescue-locator/internal/middleware"
	"emergency-rescue-locator/internal/repositories"
	"emergency-rescue-locator/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(cfg))
	router.Use(middleware.RateLimitMiddleware(cfg))

	userRepo := repositories.NewUserRepository(db)
	emergencyRepo := repositories.NewEmergencyRepository(db)
	locationRepo := repositories.NewLocationRepository(db)

	authService := services.NewAuthService(userRepo, cfg)
	emergencyService := services.NewEmergencyService(emergencyRepo, locationRepo)
	locationService := services.NewLocationService(locationRepo, emergencyRepo)

	authController := controllers.NewAuthController(authService)
	emergencyController := controllers.NewEmergencyController(emergencyService, locationService)
	locationController := controllers.NewLocationController(locationService)
	adminController := controllers.NewAdminController(emergencyService, locationService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "emergency-rescue-locator"})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.GET("/profile", middleware.AuthMiddleware(authService), authController.Profile)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			emergencies := protected.Group("/emergencies")
			{
				emergencies.POST("", emergencyController.Create)
				emergencies.GET("/active", emergencyController.GetActive)
				emergencies.GET("/:id", emergencyController.GetByID)
				emergencies.POST("/:id/cancel", emergencyController.Cancel)
				emergencies.POST("/:id/location", locationController.UpdateLocation)
				emergencies.GET("/:id/location/latest", locationController.GetLatest)
				emergencies.GET("/:id/location/history", locationController.GetHistory)
			}
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService))
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/emergencies", adminController.ListActive)
			admin.GET("/emergencies/search", adminController.Search)
			admin.GET("/emergencies/:id", adminController.GetDetails)
			admin.POST("/emergencies/:id/resolve", adminController.Resolve)
			admin.GET("/stats", adminController.Stats)
		}
	}

	return router
}
