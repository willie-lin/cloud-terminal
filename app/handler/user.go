package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CreateUser 创建一个新用户
func CreateUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		u := new(ent.User)
		if err := c.Bind(&u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
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
		users, err := client.User.Query().All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying users: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying users from database"})
		}
		return c.JSON(http.StatusOK, users)
	}
}

// GetUserByUsername 根据用户名查找
func GetUserByUsername(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// DTO
		type UserDTO struct {
			Username string `json:"username"`
		}
		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		us, err := client.User.Query().Where(user.UsernameEQ(dto.Username)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		return c.JSON(http.StatusOK, us.Username)
	}
}

// GetUserByEmail 根据邮箱查找用户
func GetUserByEmail(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// EmailDTO
		type EmailDTO struct {
			Email string `json:"email"`
		}

		dto := new(EmailDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ue, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		return c.JSON(http.StatusOK, ue)
	}
}

// UpdateUser 更新一个用户
func UpdateUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(&u); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// 检查ID是否有效
		if u.ID == uuid.Nil {
			return c.JSON(http.StatusBadRequest, "Invalid ID")
		}

		ue, err := client.User.UpdateOneID(u.ID).SetEmail(u.Email).SetNickname(u.Nickname).Save(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		if err != nil {
			log.Printf("Error updating user: %v", err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, ue)
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

// UploadFile UploadAvatar upload Avatar
func UploadFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		// 检查文件类型
		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			return errors.New("只允许上传图片文件")
		}

		// 检查文件大小
		if file.Size > 2*1024*1024 { // 限制为2MB
			return errors.New("文件太大，超过了2MB的限制")
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// 为文件重新命名
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

		// 创建目标文件
		dst, err := os.Create(filepath.Join("picture", newFileName))
		if err != nil {
			return err
		}
		defer dst.Close()

		// 将源文件复制到目标文件
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// 服务器地址
		serverAddress := "https://127.0.0.1/"
		// 返回文件的路径
		return c.String(http.StatusOK, serverAddress+filepath.Join("picture", newFileName))
	}
}

func UpdateUserInfo(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		type UpdateUserDTO struct {
			Email    string `json:"email"`
			Nickname string `json:"nickname"`
			Avatar   string `json:"avatar"`
			Phone    string `json:"phone"`
			Bio      string `json:"bio"`
		}

		dto := new(UpdateUserDTO)
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

		// 更新用户信息
		_, err = client.User.
			UpdateOne(ua).
			SetNickname(dto.Nickname).
			SetAvatar(dto.Avatar).
			SetPhone(dto.Phone).
			SetBio(dto.Bio).
			Save(context.Background())
		if err != nil {
			log.Printf("Error updating user info: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating user info in database"})
		}
		return c.JSON(http.StatusOK, map[string]string{"img": ua.Avatar})
	}
}

func UpdateUserByUUID(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(ent.User)
		if err := c.Bind(&u); err != nil {
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
