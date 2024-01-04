package api

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pquerna/otp/totp"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// CheckEmail 检查邮箱是否已经存在
func CheckEmail(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// 检查邮箱是否已经存在
		exists, err := client.User.Query().Where(user.EmailEQ(u.Email)).Exist(context.Background())
		if err != nil {
			log.Printf("Error checking email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}

func generateUsername() string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var digits = []rune("0123456789")

	var b strings.Builder
	for i := 0; i < 6; i++ {
		b.WriteRune(letters[rand.Intn(len(letters))])
	}
	for i := 0; i < 10; i++ {
		b.WriteRune(digits[rand.Intn(len(digits))])
	}
	return b.String()
}

func RegisterUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := generateUsername()
		u := new(ent.User)
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// 使用你的方法来创建密码的哈希值
		hashedPassword, err := utils.GenerateFromPassword([]byte(u.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		user, err := client.User.Create().SetEmail(u.Email).SetUsername(username).SetPassword(string(hashedPassword)).Save(context.Background())
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, user)
	}
}

func LoginUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		us, err := client.User.Query().Where(user.EmailEQ(u.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, "User-not-found")
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error-querying-user")
		}

		// 使用你的方法来验证密码和哈希值是否匹配
		if err := utils.CompareHashAndPassword([]byte(us.Password), []byte(u.Password)); err != nil {
			//return c.JSON(http.StatusForbidden, map[string]string{"error": "Invalid-password"})
			return c.JSON(http.StatusForbidden, "Invalid-password")
		}
		// 检查用户是否已经绑定了二次验证（2FA）
		if us.TotpSecret != "" {
			// 验证用户提供的OTP
			otp := u.TotpSecret
			fmt.Println(otp)
			//valid := utils.ValidateOTP(us.TotpSecret, otp)
			valid := totp.Validate(otp, us.TotpSecret)
			if !valid {
				return c.JSON(http.StatusForbidden, "Invalid-OTP")
			}
		}

		// 更新 last_login_time 字段
		us, err = client.User.
			UpdateOneID(us.ID).
			SetLastLoginTime(time.Now()).
			Save(context.Background())
		if err != nil {
			log.Printf("Error updating last login time: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error updating last login time")
		}

		// 生成JWT
		token, err := utils.GenerateToken(us.Username)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error generating token")
		}

		// 生成Refresh Token
		refreshToken, err := utils.GenerateRefreshToken(us.Username)
		if err != nil {
			log.Printf("Error generating refresh token: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error generating refresh token")
		}

		// 登录成功后，保存用户的登录信息到session
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7, // 设置session的过期时间
			HttpOnly: true,
		}
		sess.Values["username"] = user.Username // 保存用户名到session
		sess.Save(c.Request(), c.Response())

		return c.JSON(http.StatusOK, map[string]string{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

// ForgotPassword  reset password
func ForgotPassword(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		ua, err := client.User.Query().Where(user.EmailEQ(u.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, "User-not-found")
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error-querying-user")
		}
		// 使用你的方法来创建密码的哈希值
		hashedPassword, err := utils.GenerateFromPassword([]byte(u.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		_, err = client.User.
			UpdateOne(ua).SetPassword(string(hashedPassword)).Save(context.Background())
		if err != nil {
			log.Printf("Error updating password: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error updating password")
		}
		return c.JSON(http.StatusOK, "Password reset successful, please log in again.")
	}
}
