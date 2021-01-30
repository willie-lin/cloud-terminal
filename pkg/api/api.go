package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"net/http"
)

func Register(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "注册成功！！！")
	}

}
func Login(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//
		u := new(ent.User)

		// decode raw to json

		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			return err
		}

		us, err := client.User.Query().Where(user.UsernameEQ(u.Username)).Only(context.Background())
		if err != nil {
			return err
		}

		err = utils.CompareHashAndPassword([]byte(us.Password), []byte(u.Password))
		if err != utils.ErrMismatchedHashAndPassword {
			fmt.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, "登录成功！！！")

	}
}
