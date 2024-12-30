package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"timedev/db"
	"timedev/sql/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func HandleGetProfessional(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type urlParam struct {
		ReferenceKey string `param:"referencekey"`
	}

	var param urlParam
	if err := c.Bind(&param); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)

	professionalValue, err := queries.GetProfessionalInfo(ctx, param.ReferenceKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	attributeValue, err := queries.ListAttributesByProfessionalId(ctx, professionalValue.IDProfessional)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	query_deleted := c.QueryParam("deleted")
	var is_delete bool
	if query_deleted == "1" {
		is_delete = true
	}
	availabilityValue, err := queries.ListAvailabilityByProfessionalId(ctx, models.ListAvailabilityByProfessionalIdParams{
		IDProfessional: professionalValue.IDProfessional,
		Deleted:        is_delete,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	blockerValue, err := queries.ListBlockerByProfessional(ctx, models.ListBlockerByProfessionalParams{
		IDProfessional: professionalValue.IDProfessional,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalValue, "attributes": attributeValue, "availability": availabilityValue, "blocker": blockerValue})
}

func HandleCreateProfessional(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	var professionalUnit models.Professional

	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(&professionalUnit); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)

	insertedProfessional, err := queries.InsertProfessional(ctx, models.InsertProfessionalParams{
		ReferenceKey:  professionalUnit.ReferenceKey,
		Nome:          professionalUnit.Nome,
		Especialidade: professionalUnit.Especialidade,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err, "description": "Failed to insert Professional"})
	}

	return c.JSON(http.StatusOK, insertedProfessional)
}

func HandleCreateAttribute(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	type urlParam struct {
		ReferenceKey   string `param:"referencekey"`
		IDAttribute    int64  `json:"id_attribute"`
		IDProfessional int64  `json:"id_professional"`
		Attribute      string `json:"attribute"`
		Value          string `json:"value"`
	}

	var param urlParam

	if err := c.Bind(&param); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)
	tx, err := db.Begin()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to initialize a transaction"})
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	professionalUnit, err := qtx.GetProfessionalInfo(ctx, param.ReferenceKey)
	fmt.Printf("param.ReferenceKey: %v\n", param.ReferenceKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Professional not found."})
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	insertedAttribute, err := qtx.InsertAttribute(ctx, models.InsertAttributeParams{
		IDProfessional: professionalUnit.IDProfessional,
		Attribute:      param.Attribute,
		Value:          param.Value,
	})
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to insert Attribute", "description": err.Error()})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, echo.Map{"user": professionalUnit, "attributes": insertedAttribute})
}
