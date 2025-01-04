package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"timedev/db"
	"timedev/sql/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func HandleListBlocker(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close(ctx)

	type urlParams struct {
		ReferenceKey string `param:"referencekey"`
		Deleted      bool   `query:"deleted"`
	}

	var params urlParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data", "description": err})
	}

	queries := models.New(db)
	professionalUnit, err := queries.GetProfessionalInfo(ctx, params.ReferenceKey)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	blockerList, err := queries.ListBlockerByProfessional(ctx, models.ListBlockerByProfessionalParams{
		IDProfessional: professionalUnit.IDProfessional,
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
	defer db.Close(ctx)

	type urlParams struct {
		ReferenceKey string    `param:"referencekey"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		Init         time.Time `json:"init"`
		End          time.Time `json:"end"`
	}
	var params urlParams
	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data", "description": err})
	}

	queries := models.New(db)

	tx, err := db.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	qtx := queries.WithTx(tx)

	professionalUnit, err := qtx.GetProfessionalInfo(ctx, params.ReferenceKey)
	if err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	blockUnit, err := qtx.InsertBlocker(ctx, models.InsertBlockerParams{
		IDProfessional: professionalUnit.IDProfessional,
		Title:          params.Title,
		Description:    pgtype.Text{String: params.Description, Valid: true},
		InitDatetime:   pgtype.Timestamp{Time: params.Init, Valid: true},
		EndDatetime:    pgtype.Timestamp{Time: params.End, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed or Nothing to see here...", "description": err.Error()})
	}

	fmt.Printf("blockUnit.IDBlocker: %v\n", blockUnit.IDBlocker)

	slotBlocked, err := qtx.UpdateSlotSetBlocker(ctx, models.UpdateSlotSetBlockerParams{
		IDProfessional: professionalUnit.IDProfessional,
		StatusEntry:    "block",
		IDBlocker:      pgtype.Int4{Int32: blockUnit.IDBlocker, Valid: true},
		InitBlocker:    pgtype.Timestamp{Time: params.Init, Valid: true},
		EndBlocker:     pgtype.Timestamp{Time: params.End, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	tx.Commit(ctx)

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalUnit, "blocker": blockUnit, "slots_blocked": slotBlocked})
}

func HandleDeleteBlocker(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close(ctx)

	type UrlParams struct {
		ReferenceKey string `param:"referencekey"`
		IdBlocker    int32  `param:"idblocker"`
	}

	var params UrlParams

	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data", "description": err})
	}

	queries := models.New(db)

	tx, err := db.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	qtx := queries.WithTx(tx)

	professionalUnit, err := qtx.GetProfessionalInfo(ctx, params.ReferenceKey)
	if err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	blockerDeleted, err := qtx.DeleteBlockerById(ctx, params.IdBlocker)
	if err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	slotChanged, err := qtx.UpdateSlotSetBlocker(ctx, models.UpdateSlotSetBlockerParams{
		IDProfessional: professionalUnit.IDProfessional,
		StatusEntry:    "open",
		InitBlocker:    blockerDeleted.InitDatetime,
		EndBlocker:     blockerDeleted.EndDatetime,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	tx.Commit(ctx)

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalUnit, "blocker": blockerDeleted, "slots_changed": slotChanged})
}
