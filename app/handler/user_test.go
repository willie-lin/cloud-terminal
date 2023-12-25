package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/config"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// 创建一个新的 Echo 实例
	e := echo.New()

	// 创建一个新的 ent.Client 实例
	client, err := config.NewClient()
	fmt.Println(client)
	if err != nil {
		log.Fatal("opening ent client", zap.Error(err))
		return
	}
	defer client.Close()

	// 创建一个新的用户数据
	user := &ent.User{
		Username:      "testuser",
		Password:      "testpassword",
		Email:         "test@example.com",
		Nickname:      "testnickname",
		TotpSecret:    "testsecret",
		UserType:      "1",
		EnableType:    "1",
		LastLoginTime: time.Now(),
	}

	// 将用户数据转换为 JSON
	userJSON, _ := json.Marshal(user)

	// 创建一个新的 HTTP 请求
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// 创建一个 ResponseRecorder 实例来记录响应
	rec := httptest.NewRecorder()

	// 创建一个新的 Echo 上下文
	c := e.NewContext(req, rec)

	// 调用 CreateUser 函数
	if assert.NoError(t, CreateUser(client)(c)) {
		// 检查 HTTP 状态码
		assert.Equal(t, http.StatusCreated, rec.Code)

		// 检查响应中的用户数据
		var userResponse ent.User
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &userResponse)) {
			assert.Equal(t, user.Username, userResponse.Username)
			assert.Equal(t, user.Email, userResponse.Email)
			assert.Equal(t, user.Nickname, userResponse.Nickname)
			assert.Equal(t, user.TotpSecret, userResponse.TotpSecret)
			assert.Equal(t, user.UserType, userResponse.UserType)
			assert.Equal(t, user.EnableType, userResponse.EnableType)
		}
	}
}
