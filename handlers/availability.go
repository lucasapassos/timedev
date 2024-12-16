package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"timedev/db"
	"timedev/repository"
	"timedev/sql/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func HandleCreateAvailability(c echo.Context) error {
	ctx := context.Background()

	db := db.OpenDBConnection()
	defer db.Close()

	var availabilitySlot models.Availability

	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(&availabilitySlot); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	var errors_list []string

	// Validate if type Availability is in range of (0,2,3)
	if !repository.IsValidTypeAvailability(availabilitySlot.TypeAvailability) {
		errors_list = append(errors_list, "Type Availability not in (0,2,3)")
	}

	if !repository.IsValidDatetime(availabilitySlot.EndDatetime) || !repository.IsValidDatetime(availabilitySlot.InitDatetime) {
		errors_list = append(errors_list, "Error to parse End or Init Datetime")
	}

	if !repository.IsValidHour(availabilitySlot.InitHour) || !repository.IsValidHour(availabilitySlot.EndHour) {
		errors_list = append(errors_list, "Error to parse End or Init Hour")
	}

	queries := models.New(db)

	// Instanciate new transaction
	tx, err := db.Begin()
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": err, "description": "Cannot initialize db transaction"})
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	// create an author
	insertedAvailability, err := qtx.InsertAvailability(ctx, models.InsertAvailabilityParams{
		IDProfessional:   availabilitySlot.IDProfessional,
		InitDatetime:     availabilitySlot.InitDatetime,
		EndDatetime:      availabilitySlot.EndDatetime,
		InitHour:         availabilitySlot.InitHour,
		EndHour:          availabilitySlot.EndHour,
		TypeAvailability: availabilitySlot.TypeAvailability,
		WeekdayName:      availabilitySlot.WeekdayName,
		Interval:         availabilitySlot.Interval,
		PriorityEntry:    availabilitySlot.PriorityEntry,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err, "description": "Cannot insert Availability."})
	}

	slots, err := repository.ComputeSlots(insertedAvailability.InitDatetime,
		insertedAvailability.EndDatetime,
		insertedAvailability.WeekdayName,
		int(insertedAvailability.Interval),
		int(insertedAvailability.TypeAvailability),
		insertedAvailability.InitHour,
		insertedAvailability.EndHour,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err, "description": "Cannot compute slots."})
	}

	var slots_added []time.Time
	var slot_non_added []string
	for _, slot := range slots {

		value_slot_return, err := qtx.GetExistingSlot(ctx, models.GetExistingSlotParams{
			IDProfessional: insertedAvailability.IDProfessional,
			Datetime:       slot.Format("2006-01-02 15:04:05+00:00"),
			PriorityEntry:  insertedAvailability.PriorityEntry,
		})
		if value_slot_return == 0 {
			err_insert := qtx.InsertSlot(ctx, models.InsertSlotParams{
				IDProfessional: insertedAvailability.IDProfessional,
				IDAvailability: insertedAvailability.IDAvailability,
				Slot:           slot,
				WeekdayName:    insertedAvailability.WeekdayName,
				Interval:       insertedAvailability.Interval,
				PriorityEntry:  insertedAvailability.PriorityEntry,
			})
			if err_insert != nil {
				slot_non_added = append(slot_non_added, fmt.Sprint(slot.Format("2006-01-02 15:04:05+00:00"), " Failed to insert."))
			}
			slots_added = append(slots_added, slot)
		} else if err != nil {
			slot_non_added = append(slot_non_added, fmt.Sprint(slot.Format("2006-01-02 15:04:05+00:00"), " Failed to get state."))
		} else {
			slot_non_added = append(slot_non_added, fmt.Sprint(slot.Format("2006-01-02 15:04:05+00:00"), " Trying to insert in a busy slot."))
		}
	}

	if len(errors_list) > 0 {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, echo.Map{"error(s)	": errors_list})
	}
	if slots_added == nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, echo.Map{"error(s)	": "None slots and availability were added. All slots were in busy slots."})
	}

	// Commit the transaction
	tx.Commit()

	return c.JSON(http.StatusOK, echo.Map{"availability ": insertedAvailability, "slots_added": slots_added, "slots_not_added": slot_non_added})

}


func HandleGetAvailability(c echo.Context) error {
	ctx := context.Background()

	db := db.OpenDBConnection()
	defer db.Close()

	queries := models.New(db)

	availabilityIdStr := c.Param("id")
	availabilityId, err := strconv.ParseInt(availabilityIdStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to convert id in int in URL param for id availability."})
	}

	// var availabilityId custom_models.AvailabilityId
	// // Bind the incoming JSON data to the userInput struct
	// if err := c.Bind(&availabilityId); err != nil {
	// 	log.Error().Err(err).Msg("Failed to bind request data")
	// 	return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	// }

	unitAvailability, err := queries.ListAvailability(ctx, availabilityId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, err)
		}

		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, unitAvailability)
}
