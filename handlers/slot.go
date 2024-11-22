package handlers

import (
	"context"
	"log"
	"net/http"
	"time"
	"timedev/db"
	"timedev/repository"
	"timedev/sql/sqldev"

	"github.com/labstack/echo/v4"
)

var ddl string

func HandleSlots(c echo.Context) error {
	ctx := context.Background()

	db := db.OpenDBConnection()
	defer db.Close()

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := sqldev.New(db)

	layout := "2006-01-02 15:04:05"
	start, _ := time.Parse(layout, "2024-10-01 00:00:00")
	end, _ := time.Parse(layout, "2024-11-07 11:30:00")

	weekDayValids := repository.CalculateWeekdayBetween(start, end, time.Tuesday, 1)

	slots, err := repository.ComputeAgenda("16:00:00", "19:00:00", weekDayValids, 30*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	// create an author
	insertedAvailability, err := queries.InsertAvailability(ctx, sqldev.InsertAvailabilityParams{
		IDProfessional: 11123,
		InitDatetime:   "2024-10-01 00:00:00",
		EndDatetime:    "2024-11-07 11:30:00",
		InitHour:       "16:00:00",
		EndHour:        "19:00:00",
		// TypeAvailability: sql.NullInt64{Int64: 1},
		WeekdayName: "Tuesday",
		Interval:    60,
	})
	if err != nil {
		log.Println(err)
	}
	log.Println(insertedAvailability)

	return c.JSON(http.StatusOK, slots)
}
