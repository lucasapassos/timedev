package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"
	"timedev/db"
	"timedev/repository"
	"timedev/sql/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func HandleCreateSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close(ctx)

	type receivedDataStruct struct {
		ReferenceKey   string       `param:"referencekey"`
		IDAvailability int32        `json:"idavailability"`
		Slot           time.Time    `json:"slot"`
		WeekdayName    time.Weekday `json:"weekday_name"`
		Interval       int32        `json:"interval"`
		PriorityEntry  int32        `json:"priority_entry"`
		StatusEntry    string       `json:"status_entry"`
	}

	var receivedData receivedDataStruct
	if err := c.Bind(&receivedData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)
	tx, err := db.Begin(ctx)

	qtx := queries.WithTx(tx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to begin transaction", "description": err.Error()})
	}
	defer tx.Rollback(ctx)

	professionalUnit, err := qtx.GetProfessionalInfo(ctx, receivedData.ReferenceKey)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Professional does not exist"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to check professional existence", "description": err.Error()})
	}

	value_slot_return, err := qtx.GetExistingSlot(ctx, models.GetExistingSlotParams{
		IDProfessional: professionalUnit.IDProfessional,
		Slot:           pgtype.Timestamp{Time: receivedData.Slot, Valid: true},
		PriorityEntry:  receivedData.PriorityEntry,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			createdSlot, err := qtx.CreateSlot(ctx, models.CreateSlotParams{
				IDProfessional: professionalUnit.IDProfessional,
				IDAvailability: pgtype.Int4{Valid: false},
				Slot:           pgtype.Timestamp{Time: receivedData.Slot, Valid: true},
				WeekdayName:    receivedData.Slot.Weekday().String(),
				Interval:       receivedData.Interval,
				PriorityEntry:  receivedData.PriorityEntry,
				StatusEntry:    receivedData.StatusEntry,
			})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create slot", "description": err.Error()})
			}
			tx.Commit(ctx)
			return c.JSON(http.StatusCreated, createdSlot)
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get slot state", "description": err.Error()})
	}

	return c.JSON(http.StatusBadRequest, echo.Map{"error": "Slot already exists", "description": value_slot_return})
}

func HandleGetSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close(ctx)

	type urlParam struct {
		SlotId  int32 `param:"idslot"`
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
	defer db.Close(ctx)

	type SlotUnit struct {
		ReferenceKey    string    `query:"reference_key"`
		IdClinica       string    `query:"idclinica"`
		SlotInit        time.Time `query:"slot_init"`
		SlotEnd         time.Time `query:"slot_end"`
		HourInit        string    `query:"hour_init"`
		HourEnd         string    `query:"hour_end"`
		IsOpen          bool      `query:"is_open"`
		Especialidade   string    `query:"especialidade"`
		IsDeleted       bool      `query:"deleted"`
		TimingReference string    `query:"timing_reference"`
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

	var is_timing_reference bool
	if slotUnit.TimingReference != "" {
		is_timing_reference = true
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

	if is_timing_reference {
		slots, err := queries.ListSlotsSummary(ctx, models.ListSlotsSummaryParams{
			SlotInit:        pgtype.Timestamp{Time: slotUnit.SlotInit, Valid: slotUnit.SlotInit != time.Time{}},
			SlotEnd:         pgtype.Timestamp{Time: slotUnit.SlotEnd, Valid: slotUnit.SlotEnd != time.Time{}},
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
			Timing:          slotUnit.TimingReference,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed or Nothing to see here...", "description": err.Error()})
		}
		return c.JSON(http.StatusOK, slots)
	} else {
		slots, err := queries.ListSlots(ctx, models.ListSlotsParams{
			SlotInit:        pgtype.Timestamp{Time: slotUnit.SlotInit, Valid: slotUnit.SlotInit != time.Time{}},
			SlotEnd:         pgtype.Timestamp{Time: slotUnit.SlotEnd, Valid: slotUnit.SlotEnd != time.Time{}},
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
}

func HandleUpdateSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close(ctx)

	type receivedDataStruct struct {
		SlotId        int32  `param:"idslot"`
		PriorityEntry int32  `json:"priority_entry"`
		StatusEntry   string `json:"status_entry"`
		Owner         string `json:"owner"`
		IDExternal    string `json:"id_external"`
	}

	var receivedData receivedDataStruct
	if err := c.Bind(&receivedData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)
	tx, err := db.Begin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to begin transaction", "description": err.Error()})
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	slotUnit, err := qtx.GetSlotById(ctx, models.GetSlotByIdParams{IDSlot: receivedData.SlotId})
	if err != nil {
		if err == pgx.ErrNoRows {
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
		slotUnit.Owner = pgtype.Text{String: receivedData.Owner, Valid: receivedData.Owner != ""}
	}

	if receivedData.IDExternal != "" {
		slotUnit.IDExternal = pgtype.Text{String: receivedData.IDExternal, Valid: receivedData.IDExternal != ""}
	}

	updatedSlot, err := queries.UpdateSlot(ctx, models.UpdateSlotParams{
		IDSlot:        slotUnit.IDSlot,
		PriorityEntry: slotUnit.PriorityEntry,
		StatusEntry:   slotUnit.StatusEntry,
		Owner:         slotUnit.Owner,
		IDExternal:    slotUnit.IDExternal,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update slot", "description": err.Error()})
	}

	tx.Commit(ctx)
	return c.JSON(http.StatusOK, updatedSlot)
}

func HandleDeleteSlot(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close(ctx)

	type urlParam struct {
		SlotId int32 `param:"idslot"`
	}

	var param urlParam
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)
	tx, err := db.Begin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to begin transaction", "description": err.Error()})
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	slotUnit, err := qtx.GetSlotById(ctx, models.GetSlotByIdParams{IDSlot: param.SlotId})
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Slot does not exist"})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	err = qtx.DeleteSlotById(ctx, slotUnit.IDSlot)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete slot", "description": err.Error()})
	}

	tx.Commit(ctx)

	return c.JSON(http.StatusOK, echo.Map{"message": "Slot deleted"})

}
