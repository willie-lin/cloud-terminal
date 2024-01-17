package api

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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
		type UserDTO struct {
			Email string `json:"email"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		// 检查邮箱是否已经存在
		exists, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Exist(context.Background())
		if err != nil {
			log.Printf("Error checking email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking email from database"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}

// 生成唯一用户名
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

// RegisterUser 用户注册
func RegisterUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := generateUsername()
		type UserDTO struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		fmt.Println(dto.Password)

		// 使用你的方法来创建密码的哈希值
		hashedPassword, err := utils.GenerateFromPassword([]byte(dto.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		us, err := client.User.Create().SetEmail(dto.Email).SetUsername(username).SetPassword(string(hashedPassword)).Save(context.Background())
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating user in database"})
		}
		return c.JSON(http.StatusCreated, map[string]string{"userID": us.ID.String()})
	}
}

type jwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// LoginUser 用户登陆
func LoginUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type LoginDTO struct {
			Email    string  `json:"email"`
			Password string  `json:"password"`
			OTP      *string `json:"otp,omitempty"`
		}

		dto := new(LoginDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		//fmt.Println(dto.OTP)
		us, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		if len(dto.Password) == 0 {
			log.Printf("Error: password is empty")
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Password is empty"})
		}
		// 假设 us.Password 是数据库中存储的哈希值
		err = utils.CompareHashAndPassword([]byte(us.Password), []byte(dto.Password))
		if err != nil {
			log.Printf("Error comparing password: %v", err)
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Invalid password"})
		}

		// 检查用户是否已经绑定了二次验证（2FA）
		if us.TotpSecret != "" {
			// 用户已经启用了OTP，所以必须提供OTP
			if dto.OTP == nil {
				log.Printf("Error: OTP is required")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "OTP is required"})
			}
			// 验证用户提供的OTP
			valid := totp.Validate(*dto.OTP, us.TotpSecret)
			if !valid {
				log.Printf("Error: Invalid OTP")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Invalid OTP"})
			}
		}
		// update user LastLoginTime
		_, err = client.User.
			UpdateOne(us).
			SetLastLoginTime(time.Now()).
			Save(context.Background())
		if err != nil {
			log.Printf("Error updating last login time: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating last login time"})
		}
		//

		// 生成JWT
		t, err := utils.CreateToken(us.Email)
		if err != nil {
			log.Printf("Error signing token: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error signing token"})
		}

		// 将token保存在HTTP-only的cookie中，并设置相关的属性
		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = t
		cookie.Expires = time.Now().Add(1 * time.Hour)
		cookie.SameSite = http.SameSiteNoneMode
		cookie.HttpOnly = true
		cookie.Secure = true // 添加这一行
		cookie.Path = "/"
		c.SetCookie(cookie)
		// 返回成功的响应

		//// 生成RefreshToken
		//refreshToken, err := utils.GenerateRefreshToken(us.Username)
		////refreshToken, err := utils.GenerateRefreshToken(string(222))
		//if err != nil {
		//	log.Printf("Error generating refresh token: %v", err)
		//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating refresh token"})
		//}

		// 登录成功后，保存用户的登录信息到session
		sess, _ := session.Get("session", c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting session"})
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 3, // 设置session的过期时间
			HttpOnly: true,
		}
		sess.Values["username"] = us.Username // 保存用户名到session
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			log.Printf("Error saving session: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error saving session"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
		//return c.JSON(http.StatusOK, map[string]string{"token": t})
	}
}

// ResetPassword  reset password
func ResetPassword(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type UserDTO struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ua, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		// 使用你的方法来创建密码的哈希值
		hashedPassword, err := utils.GenerateFromPassword([]byte(dto.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		_, err = client.User.
			UpdateOne(ua).SetPassword(string(hashedPassword)).Save(context.Background())
		if err != nil {
			log.Printf("Error updating password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating password in database"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successful, please log in again."})
	}
}
