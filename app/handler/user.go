package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// CreateUser 创建一个新用户
func CreateUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		u := new(ent.User)

		// 直接解析raw数据为json
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		pwd, err := utils.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		u.Password = string(pwd)

		user, err := client.User.Create().
			SetUsername(u.Username).
			SetPassword(u.Password).
			SetEmail(u.Email).
			SetNickname(u.Nickname).
			SetTotpSecret(u.TotpSecret).
			SetEnableType(u.EnableType).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			Save(context.Background())
		if ent.IsConstraintError(err) {
			return c.JSON(http.StatusConflict, err.Error())
		}
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, user)
	}
}

// GetAllUsers 获取所有用户
func GetAllUsers(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//users, err := client.User.Query().All(context.Background())
		users, err := client.User.Query().All(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, users)
	}
}

// GetUserByUsername 根据用户名查找
func GetUserByUsername(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		//直接解析raw数据为json
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		user, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, user)
	}
}

// GetUserByEmail 根据邮箱查找用户
func GetUserByEmail(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		user, err := client.User.Query().Where(user.EmailEQ(u.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, user)
	}
}

// UpdateUser 更新一个用户
func UpdateUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// 检查ID是否有效
		if u.ID == uuid.Nil {
			return c.JSON(http.StatusBadRequest, "Invalid ID")
		}

		user, err := client.User.UpdateOneID(u.ID).SetEmail(u.Email).SetNickname(u.Nickname).Save(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		if err != nil {
			log.Printf("Error updating user: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, user)
	}
}

// GetUserByUUID 根据ID查找用户
func GetUserByUUID(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Printf("Invalid UUID: %v", err)
			return c.JSON(http.StatusBadRequest, "Invalid UUID")
		}

		user, err := client.User.Query().Where(user.ID(id)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, "User not found")
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error querying user")
		}
		return c.JSON(http.StatusOK, user)
	}
}

// DeleteUserByUUID  删除一个用户
func DeleteUserByUUID(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Printf("Invalid UUID: %v", err)
			return c.JSON(http.StatusBadRequest, "Invalid UUID")
		}

		err = client.User.DeleteOneID(id).Exec(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, "User not found")
		}
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error deleting user")
		}
		return c.NoContent(http.StatusNoContent)
	}
}

//
//func GetUserByUUID(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		idParam := c.Param("id")
//		id, err := uuid.Parse(idParam)
//		if err != nil {
//			log.Printf("Invalid UUID: %v", err)
//			return c.JSON(http.StatusBadRequest, "Invalid UUID")
//		}
//
//		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//		defer cancel()
//
//		user, err := client.User.Query().Where(user.ID(id)).Only(ctx)
//		if ent.IsNotFound(err) {
//			log.Printf("User not found: %v", err)
//			return c.JSON(http.StatusNotFound, "User not found")
//		}
//		if err != nil {
//			log.Printf("Error querying user: %v", err)
//			return c.JSON(http.StatusInternalServerError, "Error querying user")
//		}
//		return c.JSON(http.StatusOK, user)
//	}
//}

//func DeleteUserByUUID(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		idParam := c.Param("id")
//		id, err := uuid.Parse(idParam)
//		if err != nil {
//			log.Printf("Invalid UUID: %v", err)
//			return c.JSON(http.StatusBadRequest, "Invalid UUID")
//		}
//
//		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//		defer cancel()
//
//		err = client.User.DeleteOneID(id).Exec(ctx)
//		if ent.IsNotFound(err) {
//			log.Printf("User not found: %v", err)
//			return c.JSON(http.StatusNotFound, "User not found")
//		}
//		if err != nil {
//			log.Printf("Error deleting user: %v", err)
//			return c.JSON(http.StatusInternalServerError, "Error deleting user")
//		}
//		return c.NoContent(http.StatusNoContent)
//	}
//}

func UpdateUserByUUID(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			log.Printf("Invalid UUID: %v", err)
			return c.JSON(http.StatusBadRequest, "Invalid UUID")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		user, err := client.User.UpdateOneID(id).SetEmail(u.Email).SetNickname(u.Nickname).Save(ctx)
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, "User not found")
		}
		if err != nil {
			log.Printf("Error updating user: %v", err)
			return c.JSON(http.StatusInternalServerError, "Error updating user")
		}
		return c.JSON(http.StatusOK, user)
	}
}
