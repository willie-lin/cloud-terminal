package api

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func Register(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "注册成功！！！")
	}

}
func Login(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		log, _ := zap.NewDevelopment()
		//
		u := new(ent.User)

		// decode raw to json

		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}

		fmt.Println(u.Password)

		us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if err != nil {
			log.Fatal("The account or password you entered is incorrect.", zap.Error(err))
			return err
		}
		fmt.Println(us.Password)

		err = utils.CompareHashAndPassword([]byte(us.Password), []byte(u.Password))
		fmt.Println(err)
		//if err != utils.ErrMismatchedHashAndPassword {
		if err != nil {
			//fmt.Println(err)
			log.Fatal("The account or password you entered is incorrect.", zap.Error(err))
			//return err
		}

		return c.JSON(http.StatusOK, "登录成功！！！")

	}
}
