package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"net/http"
)

func Hello(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		//return c.JSON(http.StatusOK, map[string]string{
		//	"Echo": "hello, welcome to cloud-terminal!!!",
		//})
		return c.JSON(http.StatusOK, "hello, welcome to cloud-terminal!!!")
	}
}
