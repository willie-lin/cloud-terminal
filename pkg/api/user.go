package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/pkg/database"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"io/ioutil"
	"net/http"
	"time"
)

// 查询所有用户
func GetAllUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		client, err := database.Client()
		if err != nil {
			panic(err)
		}
		//user := new(ent.User)
		users, err := client.User.Query().All(context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, users)
	}
}

// 创建用户
func CreateUser() echo.HandlerFunc {
	//return func(c echo.Context, client *ent.Client) (*ent.User, error) {
	return func(c echo.Context) (err error) {
		//return func(c echo.Context) error {
		//var client *ent.Client
		//fmt.Println(client)
		client, err := database.Client()

		if err != nil {
			panic(err)
		}
		fmt.Println(client)
		user := new(ent.User)

		result, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			fmt.Println("接收数据失败：", err)
			return err
		}
		fmt.Println(result)
		err = json.Unmarshal(result, &user)
		if err != nil {
			fmt.Println("json解析错误", err)
			return err
		}
		fmt.Println(user.ID)
		user.ID = utils.UUID()
		fmt.Println(user.ID)

		pwd, err := utils.GenerateFromPassword([]byte(user.Password))
		if err != nil {
			fmt.Println("加密密码失败", err)
			return err
		}
		fmt.Println(pwd)
		user.Password = string(pwd)
		fmt.Println(pwd)
		u, err := client.User.Create().
			SetID(user.ID).
			SetUsername(user.Username).
			SetPassword(user.Password).
			SetNickname(user.Nickname).
			SetTotpSecret(user.TotpSecret).
			SetOnline(user.Online).
			SetEnable(user.Enable).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			SetType(user.Type).Save(context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &u)

	}

}

func DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		client, err := database.Client()
		if err != nil {
			panic(err)
			return err
		}
		u := new(ent.User)
		result, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			fmt.Println("接收数据失败：", err)
			return err
		}
		err = json.Unmarshal(result, &u)
		if err != nil {
			fmt.Println("json解析错误", err)
			return err
		}
		err = client.User.DeleteOne(u).Exec(context.Background())
		if err != nil {
			panic(err)
			fmt.Println("删除出错！")
			return err
		}
		return c.NoContent(http.StatusNoContent)
	}

}

//u, err := client.User.Create().SetID()
//var us ent.User
//if err := c.Bind(&us); err != nil {
//	return err
//}
//
//pass, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
//if err != nil {
//	fmt.Println(err)
//}
//us.Password = string(pass)
//us.ID = utils.UUID()

//u1:=User{}
//u1.Password=encodePWD //模拟从数据库中读取到的 经过bcrypt.GenerateFromPassword处理的密码值
//loginPwd:="pwd" //用户登录时输入的密码
//// 密码验证
//err = bcrypt.CompareHashAndPassword([]byte(u1.Password), []byte(loginPwd)) //验证（对比）
//if err != nil {
//	fmt.Println("pwd wrong")
//} else {
//	fmt.Println("pwd ok")
//}
