package router

import (
	"timedev/handlers"

	"github.com/labstack/echo/v4"
)

// Setup up API routes
func SetupRoutes(app *echo.Echo) {
	api := app.Group("/api")
	api.GET("/slots", handlers.HandleListSlots)
	api.POST("/slot", handlers.HandleCreateAvailability)
	api.GET("/slot/:id", handlers.HandleGetAvailability)

	// // Group that requires authentication
	// api := app.Group("/api")
	// api.Use(middleware.KeycloakJWTMiddleware)
	// api.GET("/auth-check", handlers.HandleAuthCheck)
	// api.GET("/users_alerts_subscriptions", handlers.HandleUserAlertSubscriptions)
}
