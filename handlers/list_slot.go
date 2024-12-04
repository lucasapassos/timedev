package handlers

import (
	"log"
	"net/http"
	"time"
	"timedev/db"
	"timedev/repository"

	"github.com/labstack/echo/v4"
)

var ddl string

func HandleListSlots(c echo.Context) error {

	db := db.OpenDBConnection()
	defer db.Close()

	layout := "2006-01-02 15:04:05"
	start, _ := time.Parse(layout, "2024-10-01 00:00:00")
	end, _ := time.Parse(layout, "2024-11-07 11:30:00")

	weekDayValids := repository.CalculateWeekdayBetween(start, end, time.Tuesday, 0)

	slots, err := repository.ComputeAgenda("16:00:00", "19:00:00", weekDayValids, 30*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, slots)
}
