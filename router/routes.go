package router

import (
	"timedev/handlers"

	"github.com/labstack/echo/v4"
)

// Setup up API routes
func SetupRoutes(app *echo.Echo) {
	api := app.Group("/api")

	// availability
	api.DELETE("/availability/:idavailability", handlers.HandleDeleteAvailability)
	api.POST("/availability", handlers.HandleCreateAvailability)
	api.GET("/availability/:id", handlers.HandleGetAvailability)

	// slots
	api.GET("/slot", handlers.HandleListSlots)

	// professional
	api.GET("/professional/:idprofessional", handlers.HandleGetProfessional)
	api.POST("/professional", handlers.HandleCreateProfessional)
	api.POST("/professional/attributes", handlers.HandleCreateAttribute)
	// // Group that requires authentication
	// api := app.Group("/api")
	// api.Use(middleware.KeycloakJWTMiddleware)
	// api.GET("/auth-check", handlers.HandleAuthCheck)
	// api.GET("/users_alerts_subscriptions", handlers.HandleUserAlertSubscriptions)
}
