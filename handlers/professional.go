package handlers

import (
	"context"
	"net/http"
	"timedev/db"
	"timedev/sql/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

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
