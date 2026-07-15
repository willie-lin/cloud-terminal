package handler

import (
	"github.com/labstack/echo/v5"
	"github.com/willie-lin/cloud-terminal/ent"
	"net/http"
)

func Hello(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		return c.JSON(http.StatusOK, "Hello, Welcome to Cloud-Terminal!!!")
	}
}
