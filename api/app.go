package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"timedev/config"
	"timedev/db"
	"timedev/logging"
	"timedev/middleware"
	"timedev/router"

	"github.com/labstack/echo/v4"
)

// MAIN - Setup and run
func SetupAndRunApp() error {
	log.Println("calling SetupAndRunApp()")

	log.Println("calling LoadENV()")
	// load env
	errEnv := config.LoadENV()
	if errEnv != nil {
		return errEnv
	}

	log.Println("calling SetupLogging()")
	// setup logging
	errLog := logging.SetupLogging()
	if errLog != nil {
		return errLog
	}

	log.Println("calling OpenDBConnection()")
	// Initialize the database connection
	dbConnection := db.OpenDBConnection()
	defer dbConnection.Close()

	ddl, errSchema := os.ReadFile("./sql/schema.sql")
	if errSchema != nil {
		log.Fatal(errSchema)
	}

	ctx := context.Background()
	// create tables
	if _, err := dbConnection.ExecContext(ctx, fmt.Sprintf("drop table if exists slot; drop table if exists availability; drop table if exists professional; drop table if exists attribute; %s", ddl)); err != nil {
		log.Fatal(err)
	}

	log.Println("calling Echo Instance()")
	// create Echo app -
	app := echo.New()

	// API versioning
	// v1 := app.Group("/v1")
	// router.SetupV1Routes(v1)

	// Uses API key header - 'XApiKey'
	// middleware.AddApiKeyAuth(app)

	// attach middleware
	middleware.Recover(app)
	middleware.Logger(app)

	// Use CORS - change AllowOrigins to suit
	middleware.AddCors(app)

	// setup routes
	router.SetupRoutes(app)

	// Add a rate limiter
	// middleware.RateLimiter(app)

	// Add compression
	// middleware.AddCompression(app)

	// get the server port
	port := os.Getenv("PORT")

	// Start the server
	err := app.Start(":" + port)
	if err != nil {
		return err
	}

	return nil
}
