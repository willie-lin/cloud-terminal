package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/pkg/config"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"net/http"
	"time"
)

// 查询所有用户
func GetAllUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		//client, err := database.Client()
		client, err := config.NewClient()
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

// 根据用户名查找
func FindUserByUsername() echo.HandlerFunc {
	return func(c echo.Context) error {
		//client, err := database.Client()
		client, err := config.NewClient()
		if err != nil {
			return err
		}

		u := new(ent.User)

		//// 接收raw数据
		//result, err := ioutil.ReadAll(c.Request().Body)
		//if err != nil {
		//	fmt.Println("ioutil.ReadAll err:", err)
		//	return err
		//}
		//// 解析raw为json
		//err = json.Unmarshal(result, &u)
		//if err != nil {
		//	fmt.Println("json.Unmarshal err:", err)
		//	return err
		//}

		//fmt.Println(u.Username)
		// 直接解析raw数据为json
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}
		//// or for DisallowUnknownFields() wrapped in your custom func
		//decoder := json.NewDecoder(c.Request().Body)
		//decoder.DisallowUnknownFields()
		//if err := decoder.Decode(&payload); err != nil {
		//	return err
		//}

		us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, &us)
	}
}

// 根据ID查找
func FindUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		//client, err := database.Client()
		client, err := config.NewClient()
		if err != nil {
			return err
		}

		u := new(ent.User)
		//// 接收raw数据
		//result, err := ioutil.ReadAll(c.Request().Body)
		//if err != nil {
		//	fmt.Println("ioutil.ReadAll err:", err)
		//	return err
		//}
		//// 解析raw为json
		//err = json.Unmarshal(result, &u)
		//if err != nil {
		//	fmt.Println("json.Unmarshal err:", err)
		//	return err
		//}
		//
		//fmt.Println(u.Username)

		// 直接解析raw数据为json
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}
		//// or for DisallowUnknownFields() wrapped in your custom func
		//decoder := json.NewDecoder(c.Request().Body)
		//decoder.DisallowUnknownFields()
		//if err := decoder.Decode(&payload); err != nil {
		//	return err
		//}

		us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if err != nil {
			return err
		}

		fmt.Println(us.ID)

		un, err := client.User.Query().Where(user.IDEQ(us.ID)).Only(context.Background())
		if err != nil {
			return err
		}
		fmt.Println(un)
		return c.JSON(http.StatusOK, &un)
	}
}

// 创建用户
func CreateUser() echo.HandlerFunc {
	//return func(c echo.Context, client *ent.Client) (*ent.User, error) {
	return func(c echo.Context) (err error) {
		//return func(c echo.Context) error {
		//var client *ent.Client
		//fmt.Println(client)
		//client, err := database.Client()
		client, err := config.NewClient()
		if err != nil {
			panic(err)
		}
		fmt.Println(client)
		u := new(ent.User)

		//// 接收raw数据
		//result, err := ioutil.ReadAll(c.Request().Body)
		//if err != nil {
		//	fmt.Println("ioutil.ReadAll err:", err)
		//	return err
		//}
		//
		//// 解析raw为json
		//err = json.Unmarshal(result, &u)
		//if err != nil {
		//	fmt.Println("json.Unmarshal err:", err)
		//	return err
		//}
		// 直接解析raw数据为json
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}
		//// or for DisallowUnknownFields() wrapped in your custom func
		//decoder := json.NewDecoder(c.Request().Body)
		//decoder.DisallowUnknownFields()
		//if err := decoder.Decode(&payload); err != nil {
		//	return err
		//}

		fmt.Println(u.ID)
		u.ID = utils.UUID()
		fmt.Println(u.ID)

		pwd, err := utils.GenerateFromPassword([]byte(u.Password))
		if err != nil {
			fmt.Println("加密密码失败", err)
			return err
		}
		fmt.Println(pwd)
		u.Password = string(pwd)
		fmt.Println(pwd)

		ur, err := client.User.Create().
			SetID(u.ID).
			SetUsername(u.Username).
			SetPassword(u.Password).
			SetNickname(u.Nickname).
			SetTotpSecret(u.TotpSecret).
			SetOnline(u.Online).
			SetEnable(u.Enable).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			SetType(u.Type).Save(context.Background())
		if err != nil {
			panic(err)
			return err
		}
		return c.JSON(http.StatusOK, &ur)
	}
}

// 更新用户
func UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		//client, err := database.Client()
		client, err := config.NewClient()
		if err != nil {
			panic(err)
		}
		u := new(ent.User)
		fmt.Println(client)

		//// 接收raw数据
		//result, err := ioutil.ReadAll(c.Request().Body)
		//if err != nil {
		//	fmt.Println("ioutil.ReadAll err:", err)
		//	return err
		//}
		//
		//// 解析raw为json
		//err = json.Unmarshal(result, &u)
		//if err != nil {
		//	fmt.Println("json.Unmarshal err:", err)
		//	return err
		//}

		// 直接解析raw数据为json
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}
		//// or for DisallowUnknownFields() wrapped in your custom func
		//decoder := json.NewDecoder(c.Request().Body)
		//decoder.DisallowUnknownFields()
		//if err := decoder.Decode(&payload); err != nil {
		//	return err
		//}

		//us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		//if err != nil {
		//	panic(err)
		//	return fmt.Errorf("failed querying user: %v", err)
		//}
		//_, err = us.Update().
		//	SetUsername(u.Username).
		//	SetPassword(u.Password).
		//	SetNickname(u.Nickname).
		//	SetTotpSecret(u.TotpSecret).
		//	SetOnline(u.Online).
		//	SetEnable(u.Enable).
		//	SetCreatedAt(time.Now()).
		//	SetUpdatedAt(time.Now()).
		//	SetType(u.Type).Save(context.Background())
		//if err != nil {
		//	panic(err)
		//	fmt.Println("删除出错！")
		//	return err
		//}

		ur, err := client.User.Update().
			Where(user.UsernameEQ(u.Username)).
			//SetPassword(string(utils.u.Password)).
			SetNickname(u.Nickname).
			SetTotpSecret(u.TotpSecret).
			SetOnline(u.Online).
			SetEnable(u.Enable).
			//SetCreatedAt(time.Now()).
			//SetUpdatedAt(time.Now()).
			SetType(u.Type).Save(context.Background())
		if err != nil {
			//panic(err)
			fmt.Println("update user err: ", err)
			return err
		}
		return c.JSON(http.StatusOK, &ur)
	}

}

// 更新用户 by  ID
func UpdateUserById() echo.HandlerFunc {
	return func(c echo.Context) error {

		//client, err := database.Client()
		client, err := config.NewClient()

		if err != nil {
			panic(err)
		}
		u := new(ent.User)
		fmt.Println(client)

		//// 接收raw数据
		//result, err := ioutil.ReadAll(c.Request().Body)
		//if err != nil {
		//	fmt.Println("ioutil.ReadAll err:", err)
		//	return err
		//}
		//// 解析raw为json
		//err = json.Unmarshal(result, &u)
		//if err != nil {
		//	fmt.Println("json.Unmarshal err:", err)
		//	return err
		//}

		// 直接解析raw数据为json
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}
		//// or for DisallowUnknownFields() wrapped in your custom func
		//decoder := json.NewDecoder(c.Request().Body)
		//decoder.DisallowUnknownFields()
		//if err := decoder.Decode(&payload); err != nil {
		//	return err
		//}

		us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if err != nil {
			//panic(err)
			return fmt.Errorf("failed querying user: %v", err)
		}

		ur, err := client.User.UpdateOneID(us.ID).
			SetNickname(u.Nickname).
			SetTotpSecret(u.TotpSecret).
			SetOnline(u.Online).
			SetEnable(u.Enable).
			//SetCreatedAt(time.Now()).
			//SetUpdatedAt(time.Now()).
			SetType(u.Type).Save(context.Background())
		if err != nil {
			//panic(err)
			fmt.Println("update user err: ", err)
			return err
		}

		return c.JSON(http.StatusOK, &ur)
	}

}

// 更新用户 by  ID
func TestBindJson() echo.HandlerFunc {
	return func(c echo.Context) error {

		//client, err := database.Client()
		client, err := config.NewClient()

		if err != nil {
			panic(err)
		}
		u := new(ent.User)
		fmt.Println(client)

		//// 接收raw数据
		//result, err := ioutil.ReadAll(c.Request().Body)
		//if err != nil {
		//	fmt.Println("ioutil.ReadAll err:", err)
		//	return err
		//}
		//
		//// 解析raw为json
		//err = json.Unmarshal(result, &u)
		//if err != nil {
		//	fmt.Println("json.Unmarshal err:", err)
		//	return err
		//}

		// 直接解析raw数据为json
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}
		//// or for DisallowUnknownFields() wrapped in your custom func
		//decoder := json.NewDecoder(c.Request().Body)
		//decoder.DisallowUnknownFields()
		//if err := decoder.Decode(&payload); err != nil {
		//	return err
		//}
		us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if err != nil {
			//panic(err)
			return fmt.Errorf("failed querying user: %v", err)
		}

		ur, err := client.User.UpdateOneID(us.ID).
			//SetUsername(u.Username).
			//SetPassword(u.Password).
			SetNickname(u.Nickname).
			SetTotpSecret(u.TotpSecret).
			SetOnline(u.Online).
			SetEnable(u.Enable).
			//SetCreatedAt(time.Now()).
			//SetUpdatedAt(time.Now()).
			SetType(u.Type).Save(context.Background())
		if err != nil {
			//panic(err)
			fmt.Println("update user err: ", err)
			return err
		}

		//us, err = client.User.Update().SetNickname()

		return c.JSON(http.StatusOK, &ur)

		//return c.NoContent(http.StatusNoContent)
	}

}

// 删除用户

func DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		//client, err := database.Client()
		client, err := config.NewClient()
		if err != nil {
			panic(err)
			return err
		}
		u := new(ent.User)

		//// 接收raw数据
		//result, err := ioutil.ReadAll(c.Request().Body)
		//if err != nil {
		//	fmt.Println("ioutil.ReadAll err:", err)
		//	return err
		//}
		//// 解析raw为json
		//err = json.Unmarshal(result, &u)
		//if err != nil {
		//	fmt.Println("json.Unmarshal err:", err)
		//	return err
		//}

		// 直接解析raw数据为json
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}
		//// or for DisallowUnknownFields() wrapped in your custom func
		//decoder := json.NewDecoder(c.Request().Body)
		//decoder.DisallowUnknownFields()
		//if err := decoder.Decode(&payload); err != nil {
		//	return err
		//}

		fmt.Println(1111)
		fmt.Println(u.Username)
		fmt.Println(22222)
		us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if err != nil {
			//panic(err)
			return fmt.Errorf("failed querying user: %v", err)
		}
		fmt.Println(us.ID)

		//err = client.User.DeleteOneID(u.ID).Exec(context.Background())
		err = client.User.DeleteOne(us).Exec(context.Background())
		if err != nil {
			//panic(err)
			fmt.Println("Delete user err: ", err)
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}

// 根据ID删除用户
func DeleteUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		//client, err := database.Client()
		client, err := config.NewClient()
		if err != nil {
			panic(err)
			return err
		}
		u := new(ent.User)

		//// 接收raw数据
		//result, err := ioutil.ReadAll(c.Request().Body)
		//if err != nil {
		//	fmt.Println("ioutil.ReadAll err:", err)
		//	return err
		//}
		//// 解析raw为json
		//err = json.Unmarshal(result, &u)
		//if err != nil {
		//	fmt.Println("json.Unmarshal err:", err)
		//	return err
		//}

		// 直接解析raw数据为json
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}
		//// or for DisallowUnknownFields() wrapped in your custom func
		//decoder := json.NewDecoder(c.Request().Body)
		//decoder.DisallowUnknownFields()
		//if err := decoder.Decode(&payload); err != nil {
		//	return err
		//}

		fmt.Println(1111)
		fmt.Println(u.Username)
		fmt.Println(22222)
		us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if err != nil {
			//panic(err)
			return fmt.Errorf("failed querying user: %v", err)
		}
		fmt.Println(us.ID)

		//err = client.User.DeleteOneID(u.ID).Exec(context.Background())
		err = client.User.DeleteOneID(us.ID).Exec(context.Background())
		if err != nil {
			//panic(err)
			fmt.Println("Delete user err: ", err)
			return err
		}
		//return c.NoContent(http.StatusNoContent)
		return c.NoContent(http.StatusOK)
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
