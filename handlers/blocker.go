package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"timedev/db"
	"timedev/sql/models"

	"github.com/labstack/echo/v4"
)

func HandleListBlocker(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type urlParams struct {
		IdProfessional int64 `param:"idprofessional"`
		Deleted        bool  `query:"deleted"`
	}

	var params urlParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data", "description": err})
	}

	queries := models.New(db)
	professionalUnit, err := queries.GetProfessionalInfo(ctx, params.IdProfessional)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	blockerList, err := queries.ListBlockerByProfessional(ctx, models.ListBlockerByProfessionalParams{
		IDProfessional: params.IdProfessional,
		Deleted:        params.Deleted,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if blockerList == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "No blockers found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalUnit, "blockerlist": blockerList})
}

func HandleCreateBlocker(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type urlParams struct {
		IDProfessional int64     `param:"idprofessional"`
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		Init           time.Time `json:"init"`
		End            time.Time `json:"end"`
	}
	var params urlParams
	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data", "description": err})
	}

	queries := models.New(db)

	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	qtx := queries.WithTx(tx)

	professionalUnit, err := qtx.GetProfessionalInfo(ctx, params.IDProfessional)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	blockUnit, err := qtx.InsertBlocker(ctx, models.InsertBlockerParams{
		IDProfessional: params.IDProfessional,
		Title:          params.Title,
		Description:    sql.NullString{String: params.Description, Valid: true},
		InitDatetime:   params.Init,
		EndDatetime:    params.End,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed or Nothing to see here...", "description": err.Error()})
	}

	fmt.Printf("blockUnit.IDBlocker: %v\n", blockUnit.IDBlocker)

	slotBlocked, err := qtx.UpdateSlotSetBlocker(ctx, models.UpdateSlotSetBlockerParams{
		IDProfessional: params.IDProfessional,
		StatusEntry:    "block",
		IDBlocker:      sql.NullInt64{Int64: blockUnit.IDBlocker, Valid: true},
		InitBlocker:    params.Init,
		EndBlocker:     params.End,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	tx.Commit()

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalUnit, "blocker": blockUnit, "slots_blocked": slotBlocked})
}

func HandleDeleteBlocker(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type UrlParams struct {
		IdProfessional int64 `param:"idprofessional"`
		IdBlocker      int64 `param:"idblocker"`
	}

	var params UrlParams

	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data", "description": err})
	}

	queries := models.New(db)

	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	qtx := queries.WithTx(tx)

	professionalUnit, err := qtx.GetProfessionalInfo(ctx, params.IdProfessional)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	blockerDeleted, err := qtx.DeleteBlockerById(ctx, params.IdBlocker)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	slotChanged, err := qtx.UpdateSlotSetBlocker(ctx, models.UpdateSlotSetBlockerParams{
		IDProfessional: params.IdProfessional,
		StatusEntry:    "open",
		InitBlocker:    blockerDeleted.InitDatetime,
		EndBlocker:     blockerDeleted.EndDatetime,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	tx.Commit()

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalUnit, "blocker": blockerDeleted, "slots_changed": slotChanged})
}
