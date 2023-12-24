package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"net/http"
)

func Enable2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Cloud-Terminal",
			AccountName: "123@123.com",
		})
		if err != nil {
			return err
		}
		//return c.JSON(http.StatusOK, map[string]string{
		//	"url": key.URL(),
		//})
		url := key.URL()
		var png []byte
		png, err = qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "image/png", png)
	}
}

func Validate2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		var body struct {
			Secret   string `json:"secret"`
			Passcode string `json:"passcode"`
		}
		if err := c.Bind(&body); err != nil {
			return err
		}
		valid := totp.Validate(body.Passcode, body.Secret)
		if valid {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "valid",
			})
		} else {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "invalid",
			})
		}
	}
}
