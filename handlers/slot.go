package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"
	"timedev/db"
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
		slotId  int64 `param:"idslot"`
		Deleted bool  `query:"deleted"`
	}

	var param urlParam
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)

	slotUnit, err := queries.GetSlotById(ctx, models.GetSlotByIdParams{
		IDSlot:  param.slotId,
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
		IDProfessional int64     `query:"id_professional"`
		IdClinica      string    `query:"idclinica"`
		SlotInit       time.Time `query:"slot_init"`
		SlotEnd        time.Time `query:"slot_end"`
		IsOpen         bool      `query:"is_open"`
		Especialidade  string    `query:"especialidade"`
	}

	var slotUnit SlotUnit
	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(&slotUnit); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	var is_professional bool
	if slotUnit.IDProfessional != 0 {
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
		IDProfessional:  slotUnit.IDProfessional,
		IsIdclinica:     is_idclinica,
		Idclinica:       strings.Split(slotUnit.IdClinica, ","),
		IsOpen:          slotUnit.IsOpen,
		IsEspecialidade: is_especialidade,
		Especialidade:   strings.Split(slotUnit.Especialidade, ","),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed or Nothing to see here...", "description": err.Error()})
	}

	return c.JSON(http.StatusOK, slots)
}
