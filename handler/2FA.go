package handler

import (
	"context"
	"encoding/base64"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/crypto"
	"github.com/willie-lin/cloud-terminal/viewer"
	"net/http"
)

// FaDTO 2FA 操作的请求体
type FaDTO struct {
	OTP    *string `json:"otp,omitempty"`
	Secret string  `json:"secret,omitempty"`
	// Email 仅保留给 Check2FA（公开路由，无 JWT 上下文）
	Email string `json:"email,omitempty"`
}

// currentUser 从 JWT viewer 上下文中获取当前认证用户
func currentUser(c *echo.Context, client *ent.Client) (*ent.User, error) {
	v := viewer.FromContext(c.Request().Context())
	if v == nil {
		return nil, nil
	}
	ctx := context.Background()
	ua, err := client.User.Query().Where(user.IDEQ(v.UserID.String())).Only(ctx)
	if err != nil {
		return nil, err
	}
	return ua, nil
}

// Enable2FA 为当前已认证的用户生成 TOTP 密钥与 QR 二维码
// 路由: POST /admin/enable-2fa（需要 JWT 认证）
func Enable2FA(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ua, err := currentUser(c, client)
		if ua == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Cloud-Terminal",
			AccountName: ua.Email,
		})
		if err != nil {
			log.Printf("Error generating TOTP key: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating TOTP key"})
		}

		png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
		if err != nil {
			log.Printf("Error encoding QR code: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error encoding QR code"})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"qrCode": base64.StdEncoding.EncodeToString(png),
			"secret": key.Secret(),
		})
	}
}

// Confirm2FA 验证用户扫码后的 OTP，确认后加密存储 totp_secret
// 路由: POST /admin/confirm-2FA
// Body: { "secret": "<base32>", "otp": "<6-digit>" }
func Confirm2FA(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		dto := new(FaDTO)
		if err := c.Bind(&dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ua, err := currentUser(c, client)
		if ua == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		if dto.OTP == nil || dto.Secret == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "secret and otp are required"})
		}

		// 用前端传回的 secret 验证 OTP（此时 secret 尚未入库）
		if !totp.Validate(*dto.OTP, dto.Secret) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid OTP code, please try again"})
		}

		// 信封加密后存储，不再明文写库
		encSecret, err := crypto.EncryptString(dto.Secret)
		if err != nil {
			log.Printf("Error encrypting TOTP secret: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error securing TOTP secret"})
		}

		ctx := context.Background()
		_, err = client.User.UpdateOne(ua).SetTotpSecret(encSecret).Save(ctx)
		if err != nil {
			log.Printf("Error saving TOTP secret: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error saving TOTP secret"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "2FA enabled successfully"})
	}
}

// Check2FA 检查指定 email 的用户是否已启用 2FA（公开路由，供登录页调用）
// 路由: POST /api/check-2FA
// Body: { "email": "..." }
func Check2FA(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		dto := new(FaDTO)
		if err := c.Bind(&dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ctx := context.Background()
		ua, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(ctx)
		if ent.IsNotFound(err) {
			// 不暴露用户不存在，防止枚举攻击
			return c.JSON(http.StatusOK, map[string]bool{"isConfirmed": false})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		return c.JSON(http.StatusOK, map[string]bool{"isConfirmed": ua.TotpSecret != ""})
	}
}

// Disable2FA 禁用 2FA，必须验证当前有效 OTP 后才能清除
// 路由: POST /admin/disable-2fa
// Body: { "otp": "<6-digit>" }
func Disable2FA(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		dto := new(FaDTO)
		if err := c.Bind(&dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ua, err := currentUser(c, client)
		if ua == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		if ua.TotpSecret == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "2FA is not enabled"})
		}

		if dto.OTP == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "OTP is required to disable 2FA"})
		}

		// 解密数据库中的 secret（兼容历史明文存量）
		plainSecret, err := crypto.DecryptString(ua.TotpSecret)
		if err != nil {
			log.Printf("Error decrypting TOTP secret: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error verifying OTP"})
		}

		if !totp.Validate(*dto.OTP, plainSecret) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Invalid OTP code"})
		}

		ctx := context.Background()
		_, err = client.User.UpdateOne(ua).ClearTotpSecret().Save(ctx)
		if err != nil {
			log.Printf("Error clearing TOTP secret: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error disabling 2FA"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "2FA disabled successfully"})
	}
}

// Reset2FA 重新生成 TOTP（内部调用 Enable2FA）
// 路由: POST /admin/reset-2fa
func Reset2FA(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return Enable2FA(client)(c)
	}
}
