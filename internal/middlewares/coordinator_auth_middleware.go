package middlewares

import (
	"encoding/json"
	"log/slog"

	"github.com/kamalsingh200238/ubu_management/internal/services"
	"github.com/kamalsingh200238/ubu_management/internal/utils"
	"github.com/labstack/echo/v4"
)

func CoordinatorAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Cookie("token")
		if err != nil {
			slog.Error("error in reading cookie", err)
			return utils.Redirect(c, "/login")
		}

		if token.Value == "" {
			return utils.Redirect(c, "/login")
		}

		jwtPayload, err := utils.ParseAndValidateJWT(token.Value)
		if err != nil {
			slog.Error("error in parsing jwt token", err)
			return utils.Redirect(c, "/login")
		}

		if jwtPayload.Role == utils.Coordinator {
			exist, _, err := services.CheckCoordinatorExistsByEmail(jwtPayload.Email)
			if err != nil {
				slog.Error("error in checking coordinator", err)
				return utils.Redirect(c, "/login")
			}

			if !exist {
				slog.Error("coordinator does not exist in the db", err)
				return utils.Redirect(c, "/")
			}

			jsonData, err := json.Marshal(jwtPayload)
			if err != nil {
				slog.Error("error in json marshal", err)
				return utils.Redirect(c, "/login")
			}

			c.Set("jwtPayload", jsonData)
			return next(c)
		}

		return utils.Redirect(c, "/login")
	}
}
