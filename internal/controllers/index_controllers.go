package controllers

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/validate"
	"github.com/kamalsingh200238/ubu_management/internal/services"
	"github.com/kamalsingh200238/ubu_management/internal/templates"
	"github.com/kamalsingh200238/ubu_management/internal/utils"
	"github.com/labstack/echo/v4"
)

func ShowHomePage(c echo.Context) error {
	return utils.Render(c, http.StatusOK, templates.HomePage())
}

func ShowLoginPage(c echo.Context) error {
	return utils.Render(c, http.StatusOK, templates.LoginPage(templates.LoginFormErrors{}))
}

func CooridnatorLogin(c echo.Context) error {
	formParams := struct {
		Email    string `json:"email" form:"email" validate:"required|email" message:"required:Field is required|email:Invalid emaiil address"`
		Password string `json:"password" form:"password" validate:"required" message:"required:Field is required"`
	}{}
	err := (&echo.DefaultBinder{}).BindBody(c, &formParams)
	if err != nil {
		slog.Error("error in getting params from body in cooridnator login", err)
	}

	formErrors := templates.LoginFormErrors{}
	v := validate.Struct(formParams)
	v.StopOnError = false
	if !v.Validate() {
		for field := range v.Errors.All() {
			switch field {
			case "email":
				formErrors.EmailError = v.Errors.FieldOne(field)
			case "password":
				formErrors.PasswrodError = v.Errors.FieldOne(field)
			}
		}
		return utils.Render(c, http.StatusOK, templates.LoginPage(formErrors))
	}

	exists, coordinator, err := services.CheckCoordinatorExistsByEmail(strings.ToLower(formParams.Email))
	if err != nil {
		slog.Error("error in checking if coordinator exist", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading: "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.String(http.StatusInternalServerError, "Interanl server error")
	}

	// if user doesn't exist
	if !exists {
		formErrors.EmailError = "No user found with this email"
		return utils.Render(c, http.StatusOK, templates.LoginPage(formErrors))
	}

	match, err := utils.ComparePasswordAndHash(formParams.Password, coordinator.PasswordHash)
	if err != nil {
		slog.Error("error in matching passwords", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading: "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// if user entered wrong password
	if !match {
		formErrors.PasswrodError = "Wrong password"
		return utils.Render(c, http.StatusOK, templates.LoginPage(formErrors))
	}

	// set claims for jwt token
	claims := utils.CustomJwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "your-issuer",
			Subject:   coordinator.Name,
			Audience:  jwt.ClaimStrings{"your-audience"},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 72)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		Role:     utils.Coordinator,
		Email:    coordinator.Email,
		Name:     coordinator.Name,
		PersonID: int(coordinator.ID),
	}

	token, err := utils.GenerateJwtToken(claims)
	if err != nil {
		slog.Error("error in generating jwt token", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading: "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.String(http.StatusInternalServerError, "Interanl server error")
	}

	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 72),
	})

	return utils.Redirect(c, "/coordinator")
}

func StudentLogin(c echo.Context) error {
	formParams := struct {
		Email    string `json:"email" form:"email" validate:"required|email" message:"required:Field is required|email:Invalid emaiil address"`
		Password string `json:"password" form:"password" validate:"required" message:"required:Field is required"`
	}{}
	err := (&echo.DefaultBinder{}).BindBody(c, &formParams)
	if err != nil {
		slog.Error("error in getting params from body in cooridnator login", err)
	}

	formErrors := templates.LoginFormErrors{}
	v := validate.Struct(formParams)
	v.StopOnError = false
	if !v.Validate() {
		for field := range v.Errors.All() {
			switch field {
			case "email":
				formErrors.EmailError = v.Errors.FieldOne(field)
			case "password":
				formErrors.PasswrodError = v.Errors.FieldOne(field)
			}
		}
		return utils.Render(c, http.StatusOK, templates.LoginPage(formErrors))
	}

	// check if student exists in the db
	exists, student, err := services.CheckStudentExistByEmail(strings.ToLower(formParams.Email))
	if err != nil {
		slog.Error("error in checking if student exist", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading: "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.String(http.StatusInternalServerError, "Interanl server error")
	}

	// if user doesn't exist
	if !exists {
		formErrors.EmailError = "No user found with this email"
		return utils.Render(c, http.StatusOK, templates.LoginPage(formErrors))
	}

	// check if password matched
	match, err := utils.ComparePasswordAndHash(formParams.Password, student.PasswordHash)
	if err != nil {
		slog.Error("error in matching passwords", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading: "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.String(http.StatusInternalServerError, "Interanl server error")
	}

	// if user entered wrong password
	if !match {
		formErrors.PasswrodError = "Wrong password"
		return utils.Render(c, http.StatusOK, templates.LoginPage(formErrors))
	}

	// check if the student is president
	isPresident, _, err := services.CheckIfStudentIsPresidentByID(int(student.ID))
	if err != nil {
		slog.Error("error in checking if student is president", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading: "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.String(http.StatusInternalServerError, "Interanl server error")
	}

	// decide the role of the student
	var role utils.Role
	if isPresident {
		role = utils.PresidentRole
	} else {
		role = utils.StudentRole
	}

	// set claims for jwt token
	claims := utils.CustomJwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "your-issuer",
			Subject:   student.Name,
			Audience:  jwt.ClaimStrings{"your-audience"},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 72)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		Role:     role,
		Name:     student.Name,
		Email:    student.Email,
		PersonID: int(student.ID),
	}

	token, err := utils.GenerateJwtToken(claims)
	if err != nil {
		slog.Error("error in generating token", err)
		htmx.NewResponse().AddTrigger(htmx.TriggerObject("alert", utils.AlertDetails{
			Heading: "Internal server error",
			Closable: true,
			Variant:  utils.AlertVariantDanger,
			Duration: 3000,
		})).Write(c.Response().Writer)
		return c.String(http.StatusInternalServerError, "Interanl server error")
	}

	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 72),
	})

	return utils.Redirect(c, "/student")
}
