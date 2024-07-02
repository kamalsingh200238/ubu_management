package controllers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/kamalsingh200238/ubu_management/internal/database"
	"github.com/kamalsingh200238/ubu_management/internal/middlewares"
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
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	students := []database.GetAllStudentsEnrolledInSocietyOrderByStudentNameRow{}
	if jwtPayload.Role == utils.PresidentRole {
		presidentInfoString := c.Get("presidentInfo").([]byte)
		presidentInfo := middlewares.PresidentInfo{}
		err = json.Unmarshal(presidentInfoString, &presidentInfo)
		if err != nil {
			slog.Error(
				"in unmarshal of json of president info",
				"error", err,
			)
			htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
				Heading:  "Internal server error",
				Closable: true,
				Variant:  utils.AlertVariantDanger,
				Duration: 3000,
			})).Write(c.Response().Writer)
			return c.NoContent(http.StatusInternalServerError)
		}

		students, err = database.DBQueries.GetAllStudentsEnrolledInSocietyOrderByStudentName(context.Background(), int32(presidentInfo.SocietyID))
		if err != nil {
			slog.Error(
				"in getting students enrolled in society",
				"error", err,
			)
			htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
				Heading:  "Internal server error",
				Closable: true,
				Variant:  utils.AlertVariantDanger,
				Duration: 3000,
			})).Write(c.Response().Writer)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	enrolledInSocieties, err := database.DBQueries.GetAllSocietiesStudentIsEnrolledIn(context.Background(), int32(jwtPayload.PersonID))
	if err != nil {
		slog.Error(
			"in fetching enrolled societies",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
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
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	return utils.Render(c, http.StatusOK, templates.StudentDashboard(templates.StudentDashboardParams{
		JWTPayload:                jwtPayload,
		EnrolledSocieties:         enrolledInSocieties,
		NotEnrolledSocieties:      notEnrolledSocieties,
		EnrolledStudentsInSociety: students,
	}))
}

func EnrollInSociety(c echo.Context) error {
	jwtPayloadString := c.Get("jwtPayload").([]byte)
	jwtPayload := utils.CustomJwtClaims{}
	err := json.Unmarshal(jwtPayloadString, &jwtPayload)
	if err != nil {
		slog.Error(
			"in unmarshal of json of jwt payload",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	// get the url params
	params := struct {
		StudentID int `param:"studentID"`
		SoceityID int `param:"societyID"`
	}{}
	err = (&echo.DefaultBinder{}).BindPathParams(c, &params)
	if err != nil {
		slog.Error(
			"in unmarshal of json of jwt payload",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	if jwtPayload.PersonID != params.StudentID {
		slog.Error(
			"user does not have authorization to change this resource",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Not Authorized to change this resource",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusUnauthorized)
	}

	_, err = database.DBQueries.EnrollStudentInSociety(context.Background(), database.EnrollStudentInSocietyParams{
		StudentID: int32(jwtPayload.PersonID),
		SocietyID: int32(params.SoceityID),
	})
	if err != nil {
		slog.Error(
			"in enrolling student in a society",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	students := []database.GetAllStudentsEnrolledInSocietyOrderByStudentNameRow{}
	if jwtPayload.Role == utils.PresidentRole {
		presidentInfoString := c.Get("presidentInfo").([]byte)
		presidentInfo := middlewares.PresidentInfo{}
		err = json.Unmarshal(presidentInfoString, &presidentInfo)
		if err != nil {
			slog.Error(
				"in unmarshal of json of president info",
				"error", err,
			)
			htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
				Heading:  "Internal server error",
				Closable: true,
				Variant:  utils.AlertVariantDanger,
				Duration: 3000,
			})).Write(c.Response().Writer)
			return c.NoContent(http.StatusInternalServerError)
		}

		students, err = database.DBQueries.GetAllStudentsEnrolledInSocietyOrderByStudentName(context.Background(), int32(presidentInfo.SocietyID))
		if err != nil {
			slog.Error(
				"in getting students enrolled in society",
				"error", err,
			)
			htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
				Heading:  "Internal server error",
				Closable: true,
				Variant:  utils.AlertVariantDanger,
				Duration: 3000,
			})).Write(c.Response().Writer)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	enrolledInSocieties, err := database.DBQueries.GetAllSocietiesStudentIsEnrolledIn(context.Background(), int32(jwtPayload.PersonID))
	if err != nil {
		slog.Error(
			"in fetching enrolled societies",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
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
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
		Heading:  "Enrolled in society successfully",
		Closable: true,
		Variant:  utils.AlertVariantSuccess,
		Duration: 3000,
	})).Write(c.Response().Writer)
	return utils.Render(c, http.StatusOK, templates.StudentDashboard(templates.StudentDashboardParams{
		JWTPayload:                jwtPayload,
		EnrolledSocieties:         enrolledInSocieties,
		NotEnrolledSocieties:      notEnrolledSocieties,
		EnrolledStudentsInSociety: students,
	}))
}

func LeaveSociety(c echo.Context) error {
	jwtPayloadString := c.Get("jwtPayload").([]byte)
	jwtPayload := utils.CustomJwtClaims{}
	err := json.Unmarshal(jwtPayloadString, &jwtPayload)
	if err != nil {
		slog.Error(
			"in unmarshal of json of jwt payload",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	// get the url params
	params := struct {
		StudentID int `param:"studentID"`
		SoceityID int `param:"societyID"`
	}{}
	err = (&echo.DefaultBinder{}).BindPathParams(c, &params)
	if err != nil {
		slog.Error(
			"in unmarshal of json of jwt payload",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	if jwtPayload.PersonID != params.StudentID {
		slog.Error(
			"user does not have authorization to change this resource",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Not Authorized to change this resource",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusUnauthorized)
	}

	_, err = database.DBQueries.LeaveSociety(context.Background(), database.LeaveSocietyParams{
		StudentID: int32(jwtPayload.PersonID),
		SocietyID: int32(params.SoceityID),
	})
	if err != nil {
		slog.Error(
			"in leaving society",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	students := []database.GetAllStudentsEnrolledInSocietyOrderByStudentNameRow{}
	if jwtPayload.Role == utils.PresidentRole {
		presidentInfoString := c.Get("presidentInfo").([]byte)
		presidentInfo := middlewares.PresidentInfo{}
		err = json.Unmarshal(presidentInfoString, &presidentInfo)
		if err != nil {
			slog.Error(
				"in unmarshal of json of president info",
				"error", err,
			)
			htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
				Heading:  "Internal server error",
				Closable: true,
				Variant:  utils.AlertVariantDanger,
				Duration: 3000,
			})).Write(c.Response().Writer)
			return c.NoContent(http.StatusInternalServerError)
		}

		students, err = database.DBQueries.GetAllStudentsEnrolledInSocietyOrderByStudentName(context.Background(), int32(presidentInfo.SocietyID))
		if err != nil {
			slog.Error(
				"in getting students enrolled in society",
				"error", err,
			)
			htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
				Heading:  "Internal server error",
				Closable: true,
				Variant:  utils.AlertVariantDanger,
				Duration: 3000,
			})).Write(c.Response().Writer)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	enrolledInSocieties, err := database.DBQueries.GetAllSocietiesStudentIsEnrolledIn(context.Background(), int32(jwtPayload.PersonID))
	if err != nil {
		slog.Error(
			"in fetching enrolled societies",
			"error", err,
		)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading:  "Internal server error",
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
			Heading:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
		Heading:  "Left in society successfully",
		Closable: true,
		Variant:  utils.AlertVariantSuccess,
		Duration: 3000,
	})).Write(c.Response().Writer)
	return utils.Render(c, http.StatusOK, templates.StudentDashboard(templates.StudentDashboardParams{
		JWTPayload:                jwtPayload,
		EnrolledSocieties:         enrolledInSocieties,
		NotEnrolledSocieties:      notEnrolledSocieties,
		EnrolledStudentsInSociety: students,
	}))
}
