package handler

import (
	"context"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"net/http"
)

// Enable2FA 用户设置2FA
func Enable2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//username := c.Param("username")
		username, err := GetCurrentUsername(c)
		if err != nil {
			return err
		}

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Cloud-Terminal",
			AccountName: username,
		})
		if err != nil {
			return err
		}
		url := key.URL()
		var png []byte
		png, err = qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			return err
		}

		// 从数据库中获取用户
		u, err := client.User.Query().Where(user.UsernameEQ(username)).Only(context.Background())
		if err != nil {
			return err
		}

		// 更新用户的 totp_secret 字段
		_, err = client.User.
			UpdateOne(u).
			SetTotpSecret(key.Secret()).
			Save(context.Background())
		if err != nil {
			return err
		}

		return c.Blob(http.StatusOK, "image/png", png)
	}
}

// GetCurrentUsername 获取当前登陆用户
func GetCurrentUsername(c echo.Context) (string, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return "", err
	}
	username, ok := sess.Values["username"]
	if !ok {
		return "", nil
	}

	return username.(string), nil
}

// Validate2FA 验证2FA
func Validate2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		var body struct {
			Passcode string `json:"passcode"`
		}
		if err := c.Bind(&body); err != nil {
			return err
		}

		u, err := client.User.Query().Where(user.UsernameEQ(username)).Only(context.Background())
		if err != nil {
			return err
		}

		valid := totp.Validate(body.Passcode, u.TotpSecret)
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

// Reset2FA 用户重新设置2FA
func Reset2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 调用 Enable2FA 函数来生成新的二维码
		return Enable2FA(client)(c)
	}
}
