package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"net/http"
)

func RegisterUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		user, err := client.User.Create().SetEmail(u.Email).SetPassword(string(hashedPassword)).Save(context.Background())
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

		user, err := client.User.Query().Where(user.EmailEQ(u.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, "User not found")
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error querying user")
		}

		// 使用你的方法来验证密码和哈希值是否匹配
		if err := utils.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
			return c.JSON(http.StatusUnauthorized, "Invalid password")
		}

		// 生成JWT
		token, err := utils.GenerateToken(user.Username)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error generating token")
		}

		// 生成Refresh Token
		refreshToken, err := utils.GenerateRefreshToken(user.Username)
		if err != nil {
			log.Printf("Error generating refresh token: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error generating refresh token")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

//func ForgotPassword(client *ent.Client, mailer Mailer) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		email := c.QueryParam("email")
//
//		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//		defer cancel()
//
//		user, err := client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
//		if ent.IsNotFound(err) {
//			log.Printf("User not found: %v", err)
//			return c.JSON(http.StatusNotFound, "User not found")
//		}
//		if err != nil {
//			log.Printf("Error querying user: %v", err)
//			return c.JSON(http.StatusInternalServerError, "Error querying user")
//		}
//
//		// 在实际应用中，你需要生成一个重置密码的链接，并通过电子邮件发送给用户
//		resetLink := generateResetLink(user)
//		err = mailer.SendMail(user.Email, "Reset your password", "Click this link to reset your password: "+resetLink)
//		if err != nil {
//			log.Printf("Error sending mail: %v", err)
//			return c.JSON(http.StatusInternalServerError, "Error sending mail")
//		}
//
//		return c.NoContent(http.StatusNoContent)
//	}
//}
