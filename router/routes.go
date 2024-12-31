package router

import (
	"timedev/handlers"

	"github.com/labstack/echo/v4"
)

// Setup up API routes
func SetupRoutes(app *echo.Echo) {
	api := app.Group("/api")

	// slots
	api.POST("/professional/:referencekey/slot", handlers.HandleCreateSlot)
	api.GET("/slot", handlers.HandleListSlots)
	api.GET("/slot/:idslot", handlers.HandleGetSlot)

	// professional
	api.GET("/professional/:referencekey", handlers.HandleGetProfessional)
	api.POST("/professional", handlers.HandleCreateProfessional)
	api.POST("/professional/:referencekey/attributes", handlers.HandleCreateAttribute)

	// availability
	api.DELETE("/professional/:referencekey/availability/:idavailability", handlers.HandleDeleteAvailability)
	api.POST("/professional/:referencekey/availability", handlers.HandleCreateAvailability)
	api.GET("/professional/:referencekey/availability", handlers.HandleListAvailability)
	api.GET("/professional/:referencekey/availability/:idavailability", handlers.HandleGetAvailability)

	// blockers
	api.POST("/professional/:referencekey/blocker", handlers.HandleCreateBlocker)
	api.GET("/professional/:referencekey/blocker", handlers.HandleListBlocker)
	api.DELETE("/professional/:referencekey/blocker/:idblocker", handlers.HandleDeleteBlocker)
	// Group that requires authentication
	// api := app.Group("/api")
	// api.Use(middleware.KeycloakJWTMiddleware)
	// api.GET("/auth-check", handlers.HandleAuthCheck)
	// api.GET("/users_alerts_subscriptions", handlers.HandleUserAlertSubscriptions)
}
