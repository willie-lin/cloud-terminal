package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
)

func CreateUserGroup(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		return nil
	}
}
