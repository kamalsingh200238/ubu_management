package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/gookit/validate"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kamalsingh200238/ubu_management/internal/database"
	"github.com/kamalsingh200238/ubu_management/internal/services"
	"github.com/kamalsingh200238/ubu_management/internal/templates"
	"github.com/kamalsingh200238/ubu_management/internal/utils"
	"github.com/labstack/echo/v4"
)

func ShowCoordinatorDashboard(c echo.Context) error {
	// return utils.Render(c, http.StatusOK, )
	societies, err := database.DBQueries.GetAllSocietiesWithPresidentWithStudentCount(context.Background())
	if err != nil {
		slog.Error("error in fetching societies", err)
		utils.Redirect(c, "internal-server-error")
	}
	return utils.Render(c, http.StatusOK, templates.CoordinatorDashboard(societies))
}

func ShowEditSocietyModal(c echo.Context) error {
	params := struct {
		SocietyID int `param:"id"`
	}{}
	err := (&echo.DefaultBinder{}).BindPathParams(c, &params)
	if err != nil {
		slog.Error("error in reading params", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	society, err := database.DBQueries.GetSocietyWithPresidentBySocietyId(context.Background(), int32(params.SocietyID))
	if err != nil {
		slog.Error("error in getting society from db", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	return utils.Render(c, http.StatusOK, templates.EditSocietyModal(templates.EditSocietyModalParams{
		SocietyID:                  int(society.SocietyID),
		SocietyName:                society.SocietyName,
		SocietyActive:              society.SocietyActive.Bool,
		SocietyPresidentEmail:      society.PresidentEmail,
		SocietyNameError:           "",
		SocietyPresidentEmailError: "",
	}))
}

func EditSociety(c echo.Context) error {
	// get the params out
	params := struct {
		SocietyID int `param:"id"`
	}{}
	err := (&echo.DefaultBinder{}).BindPathParams(c, &params)
	if err != nil {
		slog.Error("in getting params", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	// the form values out
	formValues := struct {
		SocietyName    string `json:"societyName" form:"societyName" validate:"required" message:"required:Field is required"`
		SocietyActive  bool   `json:"societyActive" form:"societyActive" validate:"bool"`
		PresidentEmail string `json:"presidentEmail" form:"presidentEmail" validate:"required|email" message:"required:field is required|email:Invalid email address"`
	}{}
	err = (&echo.DefaultBinder{}).BindBody(c, &formValues)
	if err != nil {
		slog.Error("in getting form values", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	societyNameError := ""
	societyPresidentEmailError := ""

	v := validate.Struct(formValues)
	v.StopOnError = false
	if !v.Validate() {
		for field := range v.Errors.All() {
			switch field {
			case "SocietyName":
				societyNameError = v.Errors.FieldOne(field)
			case "PresidentEmail":
				societyPresidentEmailError = v.Errors.FieldOne(field)
			}
		}
		htmx.NewResponse().Retarget("#modal-wrapper").Reswap("innerHTML").Reselect("#edit-society-dialog").Write(c.Response().Writer)
		return utils.Render(c, http.StatusOK, templates.EditSocietyModal(templates.EditSocietyModalParams{
			SocietyID:                  params.SocietyID,
			SocietyName:                formValues.SocietyName,
			SocietyActive:              formValues.SocietyActive,
			SocietyPresidentEmail:      formValues.PresidentEmail,
			SocietyNameError:           societyNameError,
			SocietyPresidentEmailError: societyPresidentEmailError,
		}))
	}

	exist, student, err := services.CheckStudentExistByEmail(formValues.PresidentEmail)
	if err != nil {
		slog.Error("in getting student", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	if !exist {
		societyPresidentEmailError = "Wrong email, no student found with this email"
		htmx.NewResponse().Retarget("#modal-wrapper").Reswap("innerHTML").Reselect("#edit-society-dialog").Write(c.Response().Writer)
		return utils.Render(c, http.StatusOK, templates.EditSocietyModal(templates.EditSocietyModalParams{
			SocietyID:                  params.SocietyID,
			SocietyName:                formValues.SocietyName,
			SocietyActive:              formValues.SocietyActive,
			SocietyPresidentEmail:      formValues.PresidentEmail,
			SocietyNameError:           societyNameError,
			SocietyPresidentEmailError: societyPresidentEmailError,
		}))
	}

	_, err = database.DBQueries.UpdateSociety(context.Background(), database.UpdateSocietyParams{
		ID:          int32(params.SocietyID),
		Name:        formValues.SocietyName,
		Active:      pgtype.Bool{Bool: formValues.SocietyActive, Valid: true},
		PresidentID: pgtype.Int4{Int32: student.ID, Valid: true},
	})
	if err != nil {
		slog.Error("in updating society", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	societies, err := database.DBQueries.GetAllSocietiesWithPresidentWithStudentCount(context.Background())
	if err != nil {
		slog.Error("in getting societies", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
		Message:  "Edited society successfully",
		Closable: true,
		Variant:  utils.AlertVariantSuccess,
		Duration: 3000,
	})).Write(c.Response().Writer)
	return utils.Render(c, http.StatusOK, templates.CoordinatorDashboard(societies))
}

func ShowCreateSocietyModal(c echo.Context) error {
	return utils.Render(c, http.StatusOK, templates.CreateSocietyModal(templates.CreateSocietyModalParams{
		SocietyName:                "",
		SocietyActive:              false,
		SocietyPresidentEmail:      "",
		SocietyNameError:           "",
		SocietyPresidentEmailError: "",
	}))
}

func CreateSociety(c echo.Context) error {
	// the form values out
	formValues := struct {
		SocietyName    string `json:"societyName" form:"societyName" validate:"required" message:"required:Field is required"`
		SocietyActive  bool   `json:"societyActive" form:"societyActive" validate:"bool"`
		PresidentEmail string `json:"presidentEmail" form:"presidentEmail" validate:"required|email" message:"required:field is required|email:Invalid email address"`
	}{}
	err := (&echo.DefaultBinder{}).BindBody(c, &formValues)
	if err != nil {
		slog.Error("in getting form values", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	societyNameError := ""
	societyPresidentEmailError := ""

	v := validate.Struct(formValues)
	v.StopOnError = false
	if !v.Validate() {
		for field := range v.Errors.All() {
			switch field {
			case "SocietyName":
				societyNameError = v.Errors.FieldOne(field)
			case "PresidentEmail":
				societyPresidentEmailError = v.Errors.FieldOne(field)
			}
		}
		htmx.NewResponse().Retarget("#modal-wrapper").Reswap("innerHTML").Reselect("#create-society-dialog").Write(c.Response().Writer)
		return utils.Render(c, http.StatusOK, templates.CreateSocietyModal(templates.CreateSocietyModalParams{
			SocietyName:                formValues.SocietyName,
			SocietyActive:              formValues.SocietyActive,
			SocietyPresidentEmail:      formValues.PresidentEmail,
			SocietyNameError:           societyNameError,
			SocietyPresidentEmailError: societyPresidentEmailError,
		}))
	}

	exist, student, err := services.CheckStudentExistByEmail(formValues.PresidentEmail)
	if err != nil {
		slog.Error("in getting student", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	if !exist {
		societyPresidentEmailError = "Wrong email, no student found with this email"
		htmx.NewResponse().Retarget("#modal-wrapper").Reswap("innerHTML").Reselect("#create-society-dialog").Write(c.Response().Writer)
		return utils.Render(c, http.StatusOK, templates.CreateSocietyModal(templates.CreateSocietyModalParams{
			SocietyName:                formValues.SocietyName,
			SocietyActive:              formValues.SocietyActive,
			SocietyPresidentEmail:      formValues.PresidentEmail,
			SocietyNameError:           societyNameError,
			SocietyPresidentEmailError: societyPresidentEmailError,
		}))
	}

	_, err = database.DBQueries.AddSociety(context.Background(), database.AddSocietyParams{
		Name:        formValues.SocietyName,
		Active:      pgtype.Bool{Bool: formValues.SocietyActive, Valid: true},
		PresidentID: pgtype.Int4{Int32: student.ID, Valid: true},
	})

	if err != nil {
		slog.Error("in creating society", err)
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			fmt.Println(e, e.Code, e.ColumnName, e.ColumnName, e.TableName, e.Position, e.DataTypeName, e.Detail, e.SchemaName)
			societyPresidentEmailError = "Student already a president"
			htmx.NewResponse().Retarget("#modal-wrapper").Reswap("innerHTML").Reselect("#create-society-dialog").Write(c.Response().Writer)
			return utils.Render(c, http.StatusOK, templates.CreateSocietyModal(templates.CreateSocietyModalParams{
				SocietyName:                formValues.SocietyName,
				SocietyActive:              formValues.SocietyActive,
				SocietyPresidentEmail:      formValues.PresidentEmail,
				SocietyNameError:           societyNameError,
				SocietyPresidentEmailError: societyPresidentEmailError,
			}))
		}
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	societies, err := database.DBQueries.GetAllSocietiesWithPresidentWithStudentCount(context.Background())
	if err != nil {
		slog.Error("in getting societies", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
		Message:  "Edited society successfully",
		Closable: true,
		Variant:  utils.AlertVariantSuccess,
		Duration: 3000,
	})).Write(c.Response().Writer)
	return utils.Render(c, http.StatusOK, templates.CoordinatorDashboard(societies))
}

func EnableSociety(c echo.Context) error {
	params := struct{
		SocietyID int `param:"id"`
	}{}
	err := (&echo.DefaultBinder{}).BindPathParams(c, &params)
	if err != nil {
		slog.Error("in getting params from url", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	_, err = database.DBQueries.SetSocietyActiveStatus(context.Background(), database.SetSocietyActiveStatusParams{
		Active: pgtype.Bool{Bool: true, Valid: true},
		ID: int32(params.SocietyID),
	})
	if err != nil {
		slog.Error("in setting society active status", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	societies, err := database.DBQueries.GetAllSocietiesWithPresidentWithStudentCount(context.Background())
	if err != nil {
		slog.Error("in getting societies", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Message:  "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.NoContent(http.StatusInternalServerError)
	}

	htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
		Message:  "Enabled society successfully",
		Closable: true,
		Variant:  utils.AlertVariantSuccess,
		Duration: 3000,
	})).Write(c.Response().Writer)
	return utils.Render(c, http.StatusOK, templates.CoordinatorDashboard(societies))
}
