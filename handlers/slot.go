package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"
	"timedev/db"
	"timedev/repository"
	"timedev/sql/models"

	"github.com/labstack/echo/v4"
)

func HandleCreateSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type receivedDataStruct struct {
		ReferenceKey   string        `param:"referencekey"`
		IDAvailability sql.NullInt64 `json:"idavailability"`
		Slot           time.Time     `json:"slot"`
		WeekdayName    time.Weekday  `json:"weekday_name"`
		Interval       int64         `json:"interval"`
		PriorityEntry  int64         `json:"priority_entry"`
		StatusEntry    string        `json:"status_entry"`
	}

	var receivedData receivedDataStruct
	if err := c.Bind(&receivedData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to begin transaction", "description": err.Error()})
	}
	defer tx.Rollback()

	queries := models.New(tx)

	professionalUnit, err := queries.GetProfessionalInfo(ctx, receivedData.ReferenceKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Professional does not exist"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check professional existence", "description": err.Error()})
	}

	value_slot_return, err := queries.GetExistingSlot(ctx, models.GetExistingSlotParams{
		IDProfessional: professionalUnit.IDProfessional,
		Datetime:       receivedData.Slot,
		PriorityEntry:  receivedData.PriorityEntry,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			createdSlot, err := queries.CreateSlot(ctx, models.CreateSlotParams{
				IDProfessional: professionalUnit.IDProfessional,
				IDAvailability: sql.NullInt64{Valid: false},
				Slot:           receivedData.Slot,
				WeekdayName:    receivedData.Slot.Weekday().String(),
				Interval:       receivedData.Interval,
				PriorityEntry:  receivedData.PriorityEntry,
				StatusEntry:    receivedData.StatusEntry,
			})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create slot", "description": err.Error()})
			}
			tx.Commit()
			return c.JSON(http.StatusCreated, createdSlot)
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get slot state", "description": err.Error()})
	}

	return c.JSON(http.StatusBadRequest, echo.Map{"error": "Slot already exists", "description": value_slot_return})
}

func HandleGetSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type urlParam struct {
		SlotId  int64 `param:"idslot"`
		Deleted bool  `query:"deleted"`
	}

	var param urlParam
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)

	slotUnit, err := queries.GetSlotById(ctx, models.GetSlotByIdParams{
		IDSlot:  param.SlotId,
		Deleted: param.Deleted,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, slotUnit)
}

func HandleListSlots(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type SlotUnit struct {
		ReferenceKey  string    `query:"reference_key"`
		IdClinica     string    `query:"idclinica"`
		SlotInit      time.Time `query:"slot_init"`
		SlotEnd       time.Time `query:"slot_end"`
		HourInit      string    `query:"hour_init"`
		HourEnd       string    `query:"hour_end"`
		IsOpen        bool      `query:"is_open"`
		Especialidade string    `query:"especialidade"`
		IsDeleted     bool      `query:"deleted"`
	}

	var slotUnit SlotUnit
	if err := c.Bind(&slotUnit); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	var is_hour bool
	if slotUnit.HourInit != "" && slotUnit.HourEnd != "" {
		if !repository.IsValidHour(slotUnit.HourInit) || !repository.IsValidHour(slotUnit.HourEnd) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid hour format"})
		} else {
			is_hour = true
		}
	}

	var is_professional bool
	if slotUnit.ReferenceKey != "" {
		is_professional = true
	}

	var is_idclinica bool
	if slotUnit.IdClinica != "" {
		is_idclinica = true
	}

	var is_especialidade bool
	if slotUnit.Especialidade != "" {
		is_especialidade = true
	}

	queries := models.New(db)

	slots, err := queries.ListSlots(ctx, models.ListSlotsParams{
		SlotInit:        slotUnit.SlotInit,
		SlotEnd:         slotUnit.SlotEnd,
		IsProfessional:  is_professional,
		ReferenceKey:    strings.Split(slotUnit.ReferenceKey, ","),
		IsIdclinica:     is_idclinica,
		Idclinica:       strings.Split(slotUnit.IdClinica, ","),
		IsOpen:          slotUnit.IsOpen,
		IsEspecialidade: is_especialidade,
		Especialidade:   strings.Split(slotUnit.Especialidade, ","),
		Deleted:         slotUnit.IsDeleted,
		IsHour:          is_hour,
		InitHour:        slotUnit.HourInit,
		EndHour:         slotUnit.HourEnd,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed or Nothing to see here...", "description": err.Error()})
	}

	return c.JSON(http.StatusOK, slots)
}

func HandleUpdateSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type receivedDataStruct struct {
		SlotId        int64  `param:"idslot"`
		PriorityEntry int64  `json:"priority_entry"`
		StatusEntry   string `json:"status_entry"`
		Owner         string `json:"owner"`
		ExternalID    string `json:"external_id"`
	}

	var receivedData receivedDataStruct
	if err := c.Bind(&receivedData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to begin transaction", "description": err.Error()})
	}
	defer tx.Rollback()

	queries := models.New(tx)

	slotUnit, err := queries.GetSlotById(ctx, models.GetSlotByIdParams{IDSlot: receivedData.SlotId})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Slot does not exist"})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	if receivedData.PriorityEntry != 0 {
		slotUnit.PriorityEntry = receivedData.PriorityEntry
	}

	if receivedData.StatusEntry != "" {
		slotUnit.StatusEntry = receivedData.StatusEntry
	}

	if receivedData.Owner != "" {
		slotUnit.Owner = sql.NullString{String: receivedData.Owner, Valid: receivedData.Owner != ""}
	}

	if receivedData.ExternalID != "" {
		slotUnit.ExternalID = sql.NullString{String: receivedData.ExternalID, Valid: receivedData.ExternalID != ""}
	}

	updatedSlot, err := queries.UpdateSlot(ctx, models.UpdateSlotParams{
		IDSlot:        slotUnit.IDSlot,
		PriorityEntry: slotUnit.PriorityEntry,
		StatusEntry:   slotUnit.StatusEntry,
		Owner:         slotUnit.Owner,
		ExternalID:    slotUnit.ExternalID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update slot", "description": err.Error()})
	}

	tx.Commit()
	return c.JSON(http.StatusOK, updatedSlot)
}

func HandleDeleteSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type urlParam struct {
		SlotId int64 `param:"idslot"`
	}

	var param urlParam
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to begin transaction", "description": err.Error()})
	}
	defer tx.Rollback()

	queries := models.New(tx)

	slotUnit, err := queries.GetSlotById(ctx, models.GetSlotByIdParams{IDSlot: param.SlotId})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Slot does not exist"})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	err = queries.DeleteSlotById(ctx, slotUnit.IDSlot)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete slot", "description": err.Error()})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, echo.Map{"message": "Slot deleted"})

}
