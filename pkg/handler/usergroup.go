package handler

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/group"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/user"
	"go.uber.org/zap"
	"net/http"
)

type GroupUser struct {
	UserId  string `json:"user_id,omitempty"`
	GroupId string `json:"group_id,omitempty"`
}

// 向用户组添加用户
func AddUserToGroup(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//var item UserGroup

		gu := new(GroupUser)
		//fmt.Println(gu)

		//
		log, _ := zap.NewDevelopment()
		if err := json.NewDecoder(c.Request().Body).Decode(&gu); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}
		//fmt.Println(1111111)
		//fmt.Println(gu.UserId)
		//fmt.Println(gu.GroupId)
		//fmt.Println(22222222)

		u, err := client.User.Query().Where(user.IDEQ(gu.UserId)).Only(context.Background())
		if err != nil {
			return err
		}
		g, err := client.Group.Query().Where(group.IDEQ(gu.GroupId)).Only(context.Background())
		if err != nil {
			return err
		}
		fmt.Println(u.ID, g.ID)

		_, err = client.Group.UpdateOne(g).AddUsers(u).Save(context.Background())
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}

// 从用户组删除用户
func DeleteUserFromGroup(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//var item UserGroup

		gu := new(GroupUser)
		//fmt.Println(gu)

		//
		log, _ := zap.NewDevelopment()
		if err := json.NewDecoder(c.Request().Body).Decode(&gu); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}
		//fmt.Println(1111111)
		//fmt.Println(gu.UserId)
		//fmt.Println(gu.GroupId)
		//fmt.Println(22222222)

		u, err := client.User.Query().Where(user.IDEQ(gu.UserId)).Only(context.Background())
		if err != nil {
			//log.Fatal("user not found", zap.Error(err))
			return err
		}

		g, err := client.Group.Query().Where(group.IDEQ(gu.GroupId)).Only(context.Background())
		if err != nil {
			log.Fatal("user not found", zap.Error(err))
			return err
		}

		fmt.Println(u.ID, g.ID)

		//users, err = g.QueryUsers().All(context.Background())
		//if err != nil {
		//	return err
		//}

		_ = client.Group.UpdateOne(g).RemoveUsers(u).Exec(context.Background())

		//_ = client.Group.UpdateOne(g).

		return c.NoContent(http.StatusOK)
	}
}

// 向用户组添加用户
func AddGroupToUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//var item UserGroup

		gu := new(GroupUser)
		//fmt.Println(gu)

		//
		log, _ := zap.NewDevelopment()
		if err := json.NewDecoder(c.Request().Body).Decode(&gu); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}
		//fmt.Println(1111111)
		//fmt.Println(gu.UserId)
		//fmt.Println(gu.GroupId)
		//fmt.Println(22222222)

		u, err := client.User.Query().Where(user.IDEQ(gu.UserId)).Only(context.Background())
		if err != nil {
			return err
		}
		g, err := client.Group.Query().Where(group.IDEQ(gu.GroupId)).Only(context.Background())
		if err != nil {
			return err
		}

		//client.Group.Query().Where(u.na).Only(context.Background())

		fmt.Println(u.ID, g.ID)

		//_ = client.User.UpdateOne(u).AddGroups(g).SaveX(context.Background())
		_, err = client.User.UpdateOne(u).AddGroups(g).Save(context.Background())
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

// 向用户组添加用户
func DeleteGroupFromUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		//var item UserGroup

		gu := new(GroupUser)
		//fmt.Println(gu)

		//
		log, _ := zap.NewDevelopment()
		if err := json.NewDecoder(c.Request().Body).Decode(&gu); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}
		//fmt.Println(1111111)
		//fmt.Println(gu.UserId)
		//fmt.Println(gu.GroupId)
		//fmt.Println(22222222)

		u, err := client.User.Query().Where(user.IDEQ(gu.UserId)).Only(context.Background())
		if err != nil {
			log.Fatal("user not found", zap.Error(err))
		}
		g, _ := client.Group.Query().Where(group.IDEQ(gu.GroupId)).Only(context.Background())

		fmt.Println(u.ID, g.ID)

		_ = client.User.UpdateOne(u).RemoveGroups(g).Exec(context.Background())

		//_ = client.Group.UpdateOne(g).

		return c.NoContent(http.StatusOK)
	}
}

func FindGroupWithUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		groups, err := client.Group.Query().Where(group.HasUsers()).All(context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &groups)
	}
}

func FindAllGroupWithUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		g := new(ent.Group)

		log, _ := zap.NewDevelopment()
		if err := json.NewDecoder(c.Request().Body).Decode(&g); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}

		users, err := client.Group.Query().Where(group.Name(g.Name)).QueryUsers().All(context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &users)
	}
}

func FindAllUserWithGroup(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		u := new(ent.User)

		log, _ := zap.NewDevelopment()
		if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}

		groups, err := client.User.Query().Where(user.Username(u.Username)).QueryGroups().All(context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &groups)
	}
}
