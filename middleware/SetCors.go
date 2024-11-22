package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Attach CORS to API https://echo.labstack.com/docs/middleware/cors
func AddCors(app *echo.Echo) {
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
}
