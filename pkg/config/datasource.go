package config

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/migrate"
	"go.uber.org/zap"
)

type DatabaseCfg struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	//DbPath   string      `json:"db_path"`
	Type string `json:"type"`
	//Typea     string     `json:"typea"`

}

//type Client struct {
//	Ent *ent.Client
//	CTX context.Context
//	Log *log.Logger
//}
//
//var (
//	driver  = "mysql"
//	baseUrl = "root:root1234@tcp(127.0.0.1:3306)/terminal?charset=utf8&parseTime=true"
//)
//
//func Client() (*ent.Client, error) {
//	//drv, err := sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/terminal?charset=utf8&parseTime=true")
//	drv, err := sql.Open(driver, baseUrl)
//	fmt.Println(drv)
//	//drv, err := sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/ent")
//	if err != nil {
//		return nil, err
//	}
//	db := drv.DB()
//	//db.SetMaxIdleConns(maxIdleConns)
//	//db.SetConnMaxLifetime(connMaxLifetime)
//	//db.SetMaxOpenConns(maxOpenConns)
//
//	db.SetMaxIdleConns(10)
//	db.SetMaxOpenConns(100)
//	db.SetConnMaxLifetime(time.Hour)
//	fmt.Println("ccc")
//	return ent.NewClient(ent.Driver(drv)), nil
//
//}

func NewClient() (*ent.Client, error) {
	var dfg = &DatabaseCfg{
		User:     "root",
		Password: "root1234",
		DbName:   "terminal",
		Host:     "db",
		//Host: "127.0.0.1",
		Port: 3306,
		Type: "mysql",
		//Type: "sqlite3",
	}
	var client *ent.Client
	var err error
	//drv, err := sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/terminal?charset=utf8&parseTime=true")
	switch dfg.Type {
	case "sqlite3":
		client, err = ent.Open(dfg.Type, fmt.Sprintf("file:%s?_busy_timeout=100000&_fk=1", dfg.DbName))
		if err != nil {
			return client, fmt.Errorf("failed opening connection to sqlite: %v", err)
		}
	case "mysql":
		client, err = ent.Open(dfg.Type, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			dfg.User, dfg.Password, dfg.Host, dfg.Port, dfg.DbName))
		if err != nil {
			return client, fmt.Errorf("failed opening connection to mysql: %v", err)
		}
	case "postgres", "postgresql":
		client, err = ent.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
			dfg.Host, dfg.Port, dfg.User, dfg.DbName, dfg.Password))
		if err != nil {
			return client, fmt.Errorf("failed opening connection to postgres: %v", err)
		}
	default:
		return client, fmt.Errorf("unknown database type")
	}
	return client, err
}

func AutoMigration(client *ent.Client, ctx context.Context) {
	log, _ := zap.NewDevelopment()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal("failed creating schema resources: %v", zap.Error(err))
		//log.Fatalf("failed creating schema resources: %v", err)
	}
}

func DebugMode(err error, client *ent.Client, ctx context.Context) {
	log, _ := zap.NewDevelopment()
	err = client.Debug().Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatal("failed creating schema resources: %v", zap.Error(err))
		//log.Fatalf("failed creating schema resources: %v", err)
	}
}
