package handler

import (
	"context"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"net/http"
)

// FaDTO 2FaDTO
type FaDTO struct {
	Email  string  `json:"email"`
	OTP    *string `json:"otp,omitempty"`
	Secret string  `json:"secret"`
}

func Enable2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		dto := new(FaDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Cloud-Terminal",
			AccountName: dto.Email,
		})
		if err != nil {
			log.Printf("Error generating TOTP key: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating TOTP key"})
		}
		url := key.URL()
		var png []byte
		png, err = qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			log.Printf("Error encoding QR code: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error encoding QR code"})
		}
		b64 := base64.StdEncoding.EncodeToString(png)

		return c.JSON(http.StatusOK, map[string]string{"qrCode": b64, "secret": key.Secret()})
	}
}

func Confirm2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		dto := new(FaDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		// 从数据库中获取用户
		ua, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}

		if dto.OTP == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "OTP is required"})
		}

		valid := totp.Validate(*dto.OTP, dto.Secret)
		if !valid {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid TOTP secret"})
		}

		// 确认OTP后，将密钥存储到数据库中
		_, err = client.User.
			UpdateOne(ua).
			SetTotpSecret(dto.Secret).
			Save(context.Background())
		if err != nil {
			log.Printf("Error saving TOTP secret: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error saving TOTP secret"})
		}

		return c.JSON(http.StatusOK, "2FA confirmed")
	}
}

func Check2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := privacy.DecisionContext(context.Background(), privacy.Allow)

		dto := new(FaDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		// 从数据库中获取用户
		ua, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(ctx)
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		// 检查用户是否已经绑定了二次验证（2FA）

		isConfirmed := ua.TotpSecret != ""
		return c.JSON(http.StatusOK, map[string]bool{"isConfirmed": isConfirmed})
	}
}

// Reset2FA 用户重新设置2FA
func Reset2FA(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 调用 Enable2FA 函数来生成新的二维码
		return Enable2FA(client)(c)
	}
}
