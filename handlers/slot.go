package handlers

import (
	"context"
	"net/http"
	"timedev/db"
	"timedev/repository"
	"timedev/sql/models"

	"github.com/labstack/echo/v4"
)

type SlotUnit struct {
	IDProfessional *int64 `query:"id_professional"`
	SlotInit       string `query:"slot_init"`
	SlotEnd        string `query:"slot_end"`
}

func HandleListSlots(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	slotUnit := new(SlotUnit)
	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(slotUnit); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	var is_professional bool
	if slotUnit.IDProfessional == nil {
		is_professional = true
	}

	if !(repository.IsValidDatetime(slotUnit.SlotInit)) || !(repository.IsValidDatetime(slotUnit.SlotEnd)) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Init or End slot for filter", "description": slotUnit.SlotInit})
	}

	queries := models.New(db)

	slots, err := queries.ListSlots(ctx, models.ListSlotsParams{
		SlotInit:       slotUnit.SlotInit,
		SlotEnd:        slotUnit.SlotEnd,
		IsProfessional: is_professional,
		IDProfessional: *slotUnit.IDProfessional,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed or Nothing to see here...", "description": err.Error()})
	}

	return c.JSON(http.StatusOK, slots)
}
