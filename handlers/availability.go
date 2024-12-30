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

func HandleListAvailability(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type urlParam struct {
		ReferenceKey string `param:"referencekey"`
		Deleted      bool   `query:"deleted"`
	}

	var params urlParam
	if err := c.Bind(&params); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)

	professionalUnit, err := queries.GetProfessionalInfo(ctx, params.ReferenceKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	availabilityValue, err := queries.ListAvailabilityByProfessionalId(ctx, models.ListAvailabilityByProfessionalIdParams{
		IDProfessional: professionalUnit.IDProfessional,
		Deleted:        params.Deleted,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, availabilityValue)
}

func HandleDeleteAvailability(c echo.Context) error {
	ctx := context.Background()

	db := db.OpenDBConnection()
	defer db.Close()

	type urlParam struct {
		ReferenceKey   string `param:"referencekey"`
		Deleted        bool   `query:"deleted"`
		IDAvailability int64  `param:"idavailability"`
	}

	var params urlParam
	if err := c.Bind(&params); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)

	tx, err := db.Begin()
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": "Failed to initialize a transaction"})
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	professionalUnit, err := qtx.GetProfessionalInfo(ctx, params.ReferenceKey)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	list_of_slots, err := qtx.ListSlotsByIdAvailability(ctx, sql.NullInt64{params.IDAvailability, true})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to load slots from availability", "description": err.Error()})
	}
	if list_of_slots == nil {
		return c.JSON(http.StatusNoContent, echo.Map{"error": "No slots to mark as deleted"})
	}

	var delete_slot_errors []string
	for _, slot := range list_of_slots {
		err := qtx.DeleteSlotById(ctx, slot)
		if err != nil {
			slotStr := strconv.FormatInt(slot, 10)
			delete_slot_errors = append(delete_slot_errors, fmt.Sprint("Failed to mark slot as deleted", slotStr))
		}
	}

	availabilityDeleted, err := qtx.DeleteAvailabilityById(ctx, params.IDAvailability)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to mark availability as deleted."})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalUnit, "deleted": availabilityDeleted, "description": "The slots marked as deleted", "errors": delete_slot_errors})
}

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

	type urlParam struct {
		ReferenceKey string `param:"referencekey"`
	}

	var params urlParam
	if err := c.Bind(&params); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	var errors_list []string

	// Validate if type Availability is in range of (0,2,3)
	if !repository.IsValidTypeAvailability(availabilitySlot.TypeAvailability) {
		errors_list = append(errors_list, "Type Availability not in (0,2,3)")
	}

	if !repository.IsValidHour(availabilitySlot.InitHour) || !repository.IsValidHour(availabilitySlot.EndHour) {
		errors_list = append(errors_list, "Error to parse End or Init Hour")
	}

	if len(errors_list) > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error(s)	": errors_list})
	}

	queries := models.New(db)

	// Instanciate new transaction
	tx, err := db.Begin()
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": err, "description": "Cannot initialize db transaction"})
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	professionalUnit, err := qtx.GetProfessionalInfo(ctx, params.ReferenceKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found"})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	listBlockers, err := qtx.ListBlockerByProfessional(ctx, models.ListBlockerByProfessionalParams{
		IDProfessional: professionalUnit.IDProfessional,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to capture blockers", "description": err})
	}

	// create an author
	insertedAvailability, err := qtx.InsertAvailability(ctx, models.InsertAvailabilityParams{
		IDProfessional:   professionalUnit.IDProfessional,
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

	slots, err := repository.ComputeSlots(
		insertedAvailability.InitDatetime,
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

		slotId, err := qtx.GetExistingSlot(ctx, models.GetExistingSlotParams{
			IDProfessional: professionalUnit.IDProfessional,
			Datetime:       slot,
			PriorityEntry:  insertedAvailability.PriorityEntry,
		})
		if (err != nil) && (err == sql.ErrNoRows) {

			type statusAndIdBlockerStruct struct {
				idBlocker   sql.NullInt64
				statusEntry string
			}
			var statusAndBlocker statusAndIdBlockerStruct
			statusAndBlocker.statusEntry = "open"
			statusAndBlocker.idBlocker = sql.NullInt64{Valid: false}

			for _, blockUnit := range listBlockers {
				if (blockUnit.InitDatetime.Before(slot)) && (blockUnit.EndDatetime.After(slot)) {
					statusAndBlocker.idBlocker = sql.NullInt64{Int64: blockUnit.IDBlocker, Valid: true}
					statusAndBlocker.statusEntry = "block"
				}
			}

			insertedSlot, err := qtx.InsertSlot(ctx, models.InsertSlotParams{
				IDProfessional: professionalUnit.IDProfessional,
				IDAvailability: sql.NullInt64{Int64: insertedAvailability.IDAvailability, Valid: true},
				Slot:           slot,
				WeekdayName:    insertedAvailability.WeekdayName,
				Interval:       insertedAvailability.Interval,
				PriorityEntry:  insertedAvailability.PriorityEntry,
				IsDeleted:      0,
				StatusEntry:    statusAndBlocker.statusEntry,
				IDBlocker:      statusAndBlocker.idBlocker,
			})
			if err != nil {
				slot_non_added = append(slot_non_added, fmt.Sprint(slot.Format("2006-01-02 15:04:05+00:00"), " Failed to insert."))
			}

			slots_added = append(slots_added, insertedSlot)

		} else if (err != nil) && (err != sql.ErrNoRows) {
			slot_non_added = append(slot_non_added, fmt.Sprint(slot.Format("2006-01-02 15:04:05+00:00"), " Failed to get state of slot."))
		} else {
			slot_non_added = append(slot_non_added, fmt.Sprint(slot.Format("2006-01-02 15:04:05+00:00"), " Trying to insert in a busy slot.", slotId))
		}
	}

	if slots_added == nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, echo.Map{"error(s)	": "None slots and availability were added. All slots were in busy slots."})
	}

	// Commit the transaction
	tx.Commit()

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalUnit, "availability": insertedAvailability, "slots_added": slots_added, "slots_not_added": slot_non_added})
}

func HandleGetAvailability(c echo.Context) error {
	ctx := context.Background()

	db := db.OpenDBConnection()
	defer db.Close()

	type urlParam struct {
		ReferenceKey   string `param:"referencekey"`
		Deleted        bool   `query:"deleted"`
		IDAvailability int64  `param:"idavailability"`
	}

	var params urlParam
	if err := c.Bind(&params); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)

	professionalUnit, err := queries.GetProfessionalInfo(ctx, params.ReferenceKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	unitAvailability, err := queries.ListAvailability(ctx, params.IDAvailability)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, err)
		}

		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalUnit, "availability": unitAvailability})
}
