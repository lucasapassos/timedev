package router

import (
	"timedev/handlers"

	"github.com/labstack/echo/v4"
)

// Setup up API routes
func SetupRoutes(app *echo.Echo) {
	api := app.Group("/api")
	api.GET("/slot", handlers.HandleListSlots)
	api.POST("/slot", handlers.HandleCreateAvailability)
	api.GET("/slot/:id", handlers.HandleGetAvailability)
	api.POST("/professional", handlers.HandleCreateProfessional)
	api.POST("/professional/attributes", handlers.HandleCreateAttribute)
	// // Group that requires authentication
	// api := app.Group("/api")
	// api.Use(middleware.KeycloakJWTMiddleware)
	// api.GET("/auth-check", handlers.HandleAuthCheck)
	// api.GET("/users_alerts_subscriptions", handlers.HandleUserAlertSubscriptions)
}
