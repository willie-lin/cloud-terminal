package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
)

// Handler
func enable2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//id := c.Param("id")
		//
		//// Get the user from the database
		//u, err := client.User.Get(c.Request().Context(), id)
		//if err != nil {
		//	return echo.NewHTTPError(http.StatusNotFound, "User not found")
		//}
		//
		//// Generate a new TOTP Key
		//totpKey, err := totp.Generate(totp.GenerateOpts{
		//	Issuer:      "MyApp",
		//	AccountName: u.Email, // Use the user's email as the AccountName
		//})
		//if err != nil {
		//	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate TOTP key")
		//}
		//
		//// Update the user with the new TOTP Secret
		//_, err = client.User.UpdateOneID(id).
		//	SetTotpSecret(totpKey.Secret()).
		//	Save(c.Request().Context())
		//if err != nil {
		//	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
		//}
		//
		//// Return the TOTP Key URL to display as a QR Code
		//return c.JSON(http.StatusOK, map[string]string{"url": totpKey.URL()})
		return nil
	}
}
