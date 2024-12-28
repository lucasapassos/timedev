package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
	"timedev/db"
	"timedev/sql/models"

	"github.com/labstack/echo/v4"
)

func HandleGetSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	slotIdStr := c.Param("idslot")
	slotId, err := strconv.ParseInt(slotIdStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to parse Slotid into number."})
	}

	queries := models.New(db)

	query_deleted := c.QueryParam("deleted")
	var is_delete bool
	if query_deleted == "1" {
		is_delete = true
	}

	slotUnit, err := queries.GetSlotById(ctx, models.GetSlotByIdParams{
		IDSlot:  slotId,
		Deleted: is_delete,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, slotUnit)
}

type SlotUnit struct {
	IDProfessional int64     `query:"id_professional"`
	IdClinica      string    `query:"idclinica"`
	SlotInit       time.Time `query:"slot_init"`
	SlotEnd        time.Time `query:"slot_end"`
	IsOpen         bool      `query:"is_open"`
	Especialidade  string    `query:"especialidade"`
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
