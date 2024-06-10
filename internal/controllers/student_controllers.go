package controllers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/kamalsingh200238/ubu_management/internal/database"
	"github.com/kamalsingh200238/ubu_management/internal/templates"
	"github.com/kamalsingh200238/ubu_management/internal/utils"
	"github.com/labstack/echo/v4"
)

func ShowStudentDashboard(c echo.Context) error {
	jwtPayloadString := c.Get("jwtPayload").([]byte)
	jwtPayload := utils.CustomJwtClaims{}
	err := json.Unmarshal(jwtPayloadString, &jwtPayload)
	if err != nil {
		slog.Error(
			"in unmarshal of json of jwt payload",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}
	enrolledInSocieties, err := database.DBQueries.GetAllSocietiesStudentIsEnrolledIn(context.Background(), int32(jwtPayload.PersonID))
	if err != nil {
		slog.Error(
			"in fetching enrolled societies",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	notEnrolledSocieties, err := database.DBQueries.GetAllSocietiesStudentIsNotEnrolledIn(context.Background(), int32(jwtPayload.PersonID))
	if err != nil {
		slog.Error(
			"in fetching not enrolled societies",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	return utils.Render(c, http.StatusOK, templates.StudentDashboard(templates.StudentDashboardParams{
		JWTPayload:           jwtPayload,
		EnrolledSocieties:    enrolledInSocieties,
		NotEnrolledSocieties: notEnrolledSocieties,
	}))
}
