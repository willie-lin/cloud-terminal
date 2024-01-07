package handler

import (
	"context"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"net/http"
)

func Enable2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(&u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		// 从数据库中获取用户
		ua, err := client.User.Query().Where(user.EmailEQ(u.Email)).Only(context.Background())
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Cloud-Terminal",
			AccountName: u.Email,
		})
		if err != nil {
			log.Printf("Error generating TOTP key: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		url := key.URL()
		var png []byte
		png, err = qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			log.Printf("Error encoding QR code: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		b64 := base64.StdEncoding.EncodeToString(png)
		// 存储密钥
		_, err = client.User.
			UpdateOne(ua).
			SetTotpSecret(key.Secret()).
			Save(context.Background())
		if err != nil {
			log.Printf("Error saving TOTP secret: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"qrCode": b64})
	}
}

func Confirm2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(&u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		// 从数据库中获取用户
		ua, err := client.User.Query().Where(user.EmailEQ(u.Email)).Only(context.Background())
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		otp := u.TotpSecret
		valid := totp.Validate(otp, ua.TotpSecret)
		if !valid {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid TOTP secret"})
		}

		return c.JSON(http.StatusOK, "2FA confirmed")
	}
}

func Check2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(&u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		// 从数据库中获取用户
		ua, err := client.User.Query().Where(user.EmailEQ(u.Email)).Only(context.Background())
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		// 检查用户是否已经绑定了二次验证（2FA）
		if ua.TotpSecret != "" {
			return c.JSON(http.StatusOK, map[string]bool{"isConfirmed": true})
		} else {
			return c.JSON(http.StatusOK, map[string]bool{"isConfirmed": false})
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
