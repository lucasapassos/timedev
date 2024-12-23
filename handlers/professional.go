package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"timedev/db"
	"timedev/sql/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func HandleGetProfessional(c echo.Context) error {
	ctx := context.Background()
	db := db.OpenDBConnection()
	defer db.Close()

	professionalIdStr := c.Param("idprofessional")
	professionalId, err := strconv.ParseInt(professionalIdStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to convert professional id as int"})
	}

	queries := models.New(db)

	professionalValue, err := queries.GetProfessionalInfo(ctx, professionalId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	attributeValue, err := queries.ListAttributesByProfessionalId(ctx, professionalId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	query_deleted := c.QueryParam("deleted")
	var is_delete bool
	if query_deleted == "1" {
		is_delete = true
	}
	availabilityValue, err := queries.ListAvailabilityByProfessionalId(ctx, models.ListAvailabilityByProfessionalIdParams{
		IDProfessional: professionalId,
		Deleted:        is_delete,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"professional": professionalValue, "attributes": attributeValue, "availability": availabilityValue})
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

	var attributeUnit models.Attribute

	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(&attributeUnit); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	queries := models.New(db)
	insertedAttribute, err := queries.InsertAttribute(ctx, models.InsertAttributeParams{
		IDProfessional: attributeUnit.IDProfessional,
		Attribute:      attributeUnit.Attribute,
		Value:          attributeUnit.Value,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to insert Attribute"})
	}
	return c.JSON(http.StatusOK, insertedAttribute)
}
