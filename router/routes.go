package router

import (
	"timedev/handlers"

	"github.com/labstack/echo/v4"
)

// Setup up API routes
func SetupRoutes(app *echo.Echo) {
	api := app.Group("/api")

	// slots
	api.GET("/slot", handlers.HandleListSlots)
	api.GET("/slot/:idslot", handlers.HandleGetSlot)

	// professional
	api.GET("/professional/:idprofessional", handlers.HandleGetProfessional)
	api.POST("/professional", handlers.HandleCreateProfessional)
	api.POST("/professional/:idprofessional/attributes", handlers.HandleCreateAttribute)

	// availability
	api.DELETE("/professional/:idprofessional/availability/:idavailability", handlers.HandleDeleteAvailability)
	api.POST("/professional/:idprofessional/availability", handlers.HandleCreateAvailability)
	api.GET("/professional/:idprofessional/availability", handlers.HandleListAvailability)
	api.GET("/professional/:idprofessional/availability/:id", handlers.HandleGetAvailability)

	// blockers
	api.POST("/professional/:idprofessional/blocker", handlers.HandleCreateBlocker)
	api.GET("/professional/:idprofessional/blocker", handlers.HandleListBlocker)
	api.DELETE("/professional/:idprofessional/blocker/:idblocker", handlers.HandleDeleteBlocker)
	// Group that requires authentication
	// api := app.Group("/api")
	// api.Use(middleware.KeycloakJWTMiddleware)
	// api.GET("/auth-check", handlers.HandleAuthCheck)
	// api.GET("/users_alerts_subscriptions", handlers.HandleUserAlertSubscriptions)
}
