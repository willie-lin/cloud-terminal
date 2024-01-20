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
		username := utils.GenerateUsername()

		type UserDTO struct {
			Email      string `json:"email"`
			Password   string `json:"password"`
			Online     bool   `json:"online"`
			EnableType bool   `json:"enable_type"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		// 使用你的方法来创建密码的哈希值
		hashedPassword, err := utils.GenerateFromPassword([]byte(dto.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		us, err := client.User.Create().
			SetEmail(dto.Email).
			SetUsername(username).
			SetPassword(string(hashedPassword)).
			SetOnline(dto.Online).
			SetEnableType(dto.EnableType).
			Save(context.Background())
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating user in database"})
		}
		return c.JSON(http.StatusCreated, map[string]string{"userID": us.ID.String()})
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
		// user response
		type UserResponse struct {
			Avatar   string `json:"avatar"`
			Nickname string `json:"nickname"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Bio      string `json:"bio"`
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
		// Map the user entity to the response struct
		response := &UserResponse{
			Avatar:   ue.Avatar,
			Nickname: ue.Nickname,
			Username: ue.Username,
			Email:    ue.Email,
			Phone:    ue.Phone,
			Bio:      ue.Bio,
		}
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateUser 更新一个用户
func UpdateUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type UpdateUserDTO struct {
			Username   string `json:"username"`
			Nickname   string `json:"nickname"`
			Phone      string `json:"phone"`
			Bio        string `json:"bio"`
			Online     bool   `json:"online"`
			EnableType bool   `json:"enable_type"`
		}

		dto := new(UpdateUserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		fmt.Println(1111111)
		fmt.Println(dto.Online)
		fmt.Println(dto.EnableType)
		fmt.Println(2222222)

		// 从数据库中获取用户
		ua, err := client.User.Query().Where(user.UsernameEQ(dto.Username)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}

		// 更新数据库
		_, err = client.User.UpdateOne(ua).
			SetNickname(dto.Nickname).
			SetPhone(dto.Phone).
			SetBio(dto.Bio).
			SetOnline(dto.Online).
			SetEnableType(dto.EnableType).
			SetLastLoginTime(time.Now()).
			Save(context.Background())
		if err != nil {
			log.Printf("Error updating user info: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating user info in database"})
		}
		//return c.JSON(http.StatusOK, map[string]string{"img": ua.Avatar})
		return c.JSON(http.StatusOK, map[string]string{"message": "User update successful"})
	}
}

// GetUserByUUID 根据ID查找用户

// DeleteUserByUsername 删除一个用户
func DeleteUserByUsername(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// UserDTO
		type UserDTO struct {
			Username string `json:"username"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ue, err := client.User.Query().Where(user.UsernameEQ(dto.Username)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}

		err = client.User.DeleteOne(ue).Exec(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting user"})
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

// UpdateUserInfo Update user info
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
		//return c.JSON(http.StatusOK, map[string]string{"img": ua.Avatar})
		return c.JSON(http.StatusOK, map[string]string{"message": "User update successful"})
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
