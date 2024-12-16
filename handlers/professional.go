package handlers

import (
	"context"
	"fmt"
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

	attributeList := &[]models.Attribute{}

	// Bind the incoming JSON data to the userInput struct
	if err := c.Bind(attributeList); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request data"})
	}

	var errorsLog []string
	queries := models.New(db)
	for _, attribute := range *attributeList {
		err := queries.InsertAttribute(ctx, models.InsertAttributeParams{
			attribute.IDProfessional,
			attribute.Attribute,
			attribute.Value,
		})
		if err != nil {
			errorsLog = append(errorsLog, fmt.Sprintf("Failed to add attribute %s", attribute.Attribute))
		}
	}

	if len(errorsLog) > 0 {
		return c.JSON(echo.ErrBadRequest.Code, errorsLog)
	}

	return c.JSON(http.StatusOK, attributeList)
}
