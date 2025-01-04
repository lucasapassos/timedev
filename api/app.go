package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
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

	location, err_tz := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err_tz != nil {
		panic(err_tz)
	}
	time.Local = location

	log.Println("calling OpenDBConnection()")

	initialize_script_db := `
  
  drop table if exists slot; 
  drop table if exists availability; 
  drop table if exists blocker;
  drop table if exists attribute;
  drop table if exists professional;  `

	ctx := context.Background()
	// Initialize the database connection
	dbConnection := db.OpenDBConnection()
	defer dbConnection.Close(ctx)

	ddl, errSchema := os.ReadFile("./sql/schema.sql")
	if errSchema != nil {
		log.Fatal(errSchema)
	}
	// create tables
	if _, err := dbConnection.Exec(ctx, fmt.Sprintf("%s; %s", initialize_script_db, ddl)); err != nil {
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
